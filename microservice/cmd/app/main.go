package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	dataMng "github.com/john98nf/UltimateMicroservice/cmd/app/dataManagement"
	"github.com/john98nf/UltimateMicroservice/cmd/app/model"
)

var mdlCtrl *dataMng.MiddlewareController

func authenticationWrapper(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		header, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"Status": "Authentication required"})
			return
		}

		// Format Bearer <token>
		parts := strings.Split(header[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"Status": "Authentication required"})
			return
		}
		tokenStr := parts[1]

		if err := mdlCtrl.ValidateToken(tokenStr); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"Status": "Authentication required"})
			return
		}

		next(c)
	}
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/signup", func(c *gin.Context) {
		var creds model.Credentials
		creds.Username = c.PostForm("Username")
		creds.Password = c.PostForm("Password")
		if creds.Username == "" || creds.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid payload"})
			return
		}

		if err := mdlCtrl.CreateUser(creds); err != nil {
			if errors.Is(err, dataMng.DuplicateResource) {
				// Typically this is not a good approach
				// Platforms must not provide the outside world
				// with info like the existance of a username
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Username taken"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "User Created"})
	})

	r.POST("/signin", func(c *gin.Context) {
		var creds model.Credentials
		creds.Username = c.PostForm("username")
		creds.Password = c.PostForm("password")
		if creds.Username == "" || creds.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Status": "Authentication Failed"})
			return
		}

		tokenString, err := mdlCtrl.GenerateTokenJWT(&creds)
		if err != nil {
			if errors.Is(err, dataMng.UserAuthenticationFailed) {
				c.JSON(http.StatusUnauthorized, gin.H{"Status": "Authentication Failed"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "User Authorised", "token": tokenString})
	})

	r.GET("/company/:id", func(c *gin.Context) {
		id, errId := uuid.Parse(c.Params.ByName("id"))
		if errId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}
		company, err := mdlCtrl.GetCompany(id)
		if err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.ResourceNotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"Status": "Resource Not Found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}
		c.JSON(http.StatusOK, company)
	})

	r.PATCH("/company/:id", authenticationWrapper(func(c *gin.Context) {
		id, err := uuid.Parse(c.Params.ByName("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}

		noNewValue := true
		var (
			name               *string
			description        *string
			employees          *int
			registrationStatus *bool
			legalType          *string
		)

		nameForm := c.PostForm("Name")
		if nameForm != "" {
			name = &nameForm
			noNewValue = false
		}
		descriptionForm := c.PostForm("Description")
		if descriptionForm != "" {
			description = &descriptionForm
			noNewValue = false
		}
		var employeesForm int
		if emp := c.PostForm("Employees"); emp != "" {
			employeesForm, err = strconv.Atoi(emp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Employees Value"})
				return
			}
			employees = &employeesForm
			noNewValue = false
		}
		var registrationStatusForm bool
		if regStatus := c.PostForm("RegistrationStatus"); regStatus != "" {
			registrationStatusForm, err = strconv.ParseBool(regStatus)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Registration Status"})
				return
			}
			registrationStatus = &registrationStatusForm
			noNewValue = false
		}
		legalTypeForm := c.PostForm("LegalType")
		if legalTypeForm != "" {
			if ok := model.VerifyCompanyType(legalTypeForm); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company type"})
				return
			}
			legalType = &legalTypeForm
			noNewValue = false
		}

		if noNewValue {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Attributes"})
			return
		}

		company, err := mdlCtrl.ModifyCompany(id,
			name,
			description,
			employees,
			registrationStatus,
			legalType)
		if err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.ResourceNotFoundError) {
				c.JSON(http.StatusConflict, gin.H{"Status": "Resource Not Found"})
			} else if errors.Is(err, dataMng.DuplicateResource) {
				c.JSON(http.StatusConflict, gin.H{"Status": "Conflict with another resource"})
			} else if errors.Is(err, dataMng.NoResourceModification) {
				c.JSON(http.StatusOK, gin.H{"Status": "Company wasn't modified"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}
		c.JSON(http.StatusOK, company)
	}))

	r.POST("/company", authenticationWrapper(func(c *gin.Context) {
		id := uuid.Nil // UUID construction is handled by the MiddlewareController
		name := c.PostForm("Name")
		description := c.PostForm("Description")
		employees, errEmployees := strconv.Atoi(c.PostForm("Employees"))
		registrationStatus, _ := strconv.ParseBool(c.PostForm("RegistrationStatus"))
		legalType := c.PostForm("LegalType")
		if errEmployees != nil || !model.VerifyCompanyType(legalType) {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Attributes"})
			return
		}

		newCompany := &model.Company{
			Id:                 id,
			Name:               name,
			Description:        description,
			Employees:          employees,
			RegistrationStatus: registrationStatus,
			LegalType:          legalType,
		}
		if err := mdlCtrl.CreateCompany(newCompany); err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.DuplicateResource) {
				c.JSON(http.StatusConflict, gin.H{"Status": "Company already exists"})
			} else if errors.Is(err, dataMng.UnavailableUUIDGeneration) {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Company creation unavailable"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}
		c.JSON(http.StatusOK, newCompany)
	}))

	r.DELETE("/company/:id", authenticationWrapper(func(c *gin.Context) {
		id, errId := uuid.Parse(c.Params.ByName("id"))
		if errId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}
		if err := mdlCtrl.DeleteCompany(id); err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.ResourceNotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"Status": "Resource Not Found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"Status": "Company deleted"})
	}))

	return r
}

func main() {

	mdlCtrl = dataMng.InitiallizeNewMiddlewareController()

	r := setupRouter()
	r.Run(":8080")
}
