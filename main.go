package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

// func for connecting to the database
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Static("/documents/temp", "./documents/temp")
	router.GET("/leaves/:id", getId)
	router.GET("/leaves", getLeaves)
	router.POST("/leave", postLeave)
	router.GET("/file/:id", getFile) // endpoint for file viewing

	router.Run(":" + port)

   
}




