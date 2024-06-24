package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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

	cfg := mysql.Config{
		User:   envFile["DBUSER"],
		Passwd: envFile["DBPASSWORD"],
		Net:    "tcp",
		Addr:   envFile["DBENDPOINT"],
		DBName: envFile["DBSCHEMA"],
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("Couldn't ping the database.")
		log.Fatal(pingErr)
	}
	log.Println("Connected to the Database!")

	r := setupRouter()
	r.Run(":8080")
}
