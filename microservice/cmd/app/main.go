package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		sId, errId := strconv.Atoi(c.Params.ByName("id"))
		if errId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}
		id := uint(sId)
		company, err := mdlCtrl.GetCompany(id)
		if err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.ResourceNotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"Status": "Resource Not Found"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Request"})
			}
			return
		}
		c.JSON(http.StatusOK, company)
	})

	r.POST("/company", func(c *gin.Context) {
		id, errId := strconv.Atoi(c.PostForm("id"))
		name := c.PostForm("name")
		description := c.PostForm("description")
		employees, errEmployees := strconv.Atoi(c.PostForm("employees"))
		registrationStatus, _ := strconv.ParseBool(c.PostForm("registrationStatus"))
		legalType := c.PostForm("legalType")
		if errId != nil || errEmployees != nil || !model.VerifyCompanyType(legalType) {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Attributes"})
			return
		}

		newCompany := &model.Company{
			Id:                 uint(id),
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
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Request"})
			return
		}
		c.JSON(http.StatusOK, newCompany)
	})

	r.DELETE("/company/:id", func(c *gin.Context) {
		sId, errId := strconv.Atoi(c.Params.ByName("id"))
		if errId != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Company Id"})
			return
		}
		id := uint(sId)
		if err := mdlCtrl.DeleteCompany(id); err != nil {
			log.Println(err)
			if errors.Is(err, dataMng.ResourceNotFoundError) {
				c.JSON(http.StatusNotFound, gin.H{"Status": "Resource Not Found"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Status": "Invalid Request"})
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
