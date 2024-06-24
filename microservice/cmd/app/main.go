package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dataMng "github.com/john98nf/UltimateMicroservice/cmd/app/dataManagement"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	// HealthCheck endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {

	envFile, err := godotenv.Read("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	var mdlCtrl *dataMng.MiddlewareController = dataMng.InitiallizeNewMiddlewareController(envFile)
	fmt.Println(mdlCtrl)

	r := setupRouter()
	r.Run(":8080")
}
