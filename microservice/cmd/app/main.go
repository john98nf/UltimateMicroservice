package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	dataMng "github.com/john98nf/UltimateMicroservice/cmd/app/dataManagement"
	"github.com/john98nf/UltimateMicroservice/cmd/app/model"
	"github.com/joho/godotenv"
)

var mdlCtrl *dataMng.MiddlewareController

func setupRouter() *gin.Engine {

	r := gin.Default()

	// HealthCheck endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// GET company
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

	r.PATCH("/company/:id", func(c *gin.Context) {
		id, errId := uuid.Parse(c.Params.ByName("id"))
		if errId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}

		// Retrieve current company
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
		noNewValue := true
		if name := c.PostForm("Name"); name != "" {
			company.Name = name
			noNewValue = false
		}
		if description := c.PostForm("Description"); description != "" {
			company.Description = description
			noNewValue = false
		}
		if emp := c.PostForm("Employees"); emp != "" {
			employees, err := strconv.Atoi(emp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Employees Value"})
				return
			}
			company.Employees = employees
			noNewValue = false
		}
		if regStatus := c.PostForm("RegistrationStatus"); regStatus != "" {
			registrationStatus, err := strconv.ParseBool(regStatus)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Registration Status"})
				return
			}
			company.RegistrationStatus = registrationStatus
			noNewValue = false
		}
		if legalType := c.PostForm("LegalType"); legalType != "" {
			if ok := model.VerifyCompanyType(legalType); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company type"})
				return
			}
			company.LegalType = legalType
			noNewValue = false
		}

		if noNewValue {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Attributes"})
			return
		}

		if err := mdlCtrl.ModifyCompany(company); err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.DuplicateResource) {
				c.JSON(http.StatusConflict, gin.H{"Status": "Conflict with another resource"})
			} else if errors.Is(err, dataMng.NoResourceModification) {
				c.JSON(http.StatusOK, gin.H{"Status": "Company wasn't modified"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
			}
			return
		}
		c.JSON(http.StatusOK, company)
	})

	r.POST("/company", func(c *gin.Context) {
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
	})

	r.DELETE("/company/:id", func(c *gin.Context) {
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
	})

	return r
}

func main() {

	envFile, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	mdlCtrl = dataMng.InitiallizeNewMiddlewareController(envFile)

	r := setupRouter()
	r.Run(":8080")
}
