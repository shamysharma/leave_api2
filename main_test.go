package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func abcMain() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	//port := os.Getenv("PORT")

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetId(t *testing.T) {
	// Set up a test database
	abcMain()

	// Create a Gin router for testing
	router := gin.Default()
	router.GET("/leaves/:id", getId)

	// Create a test request
	req, err := http.NewRequest("GET", "/leaves/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Serve the request and record the response
	router.ServeHTTP(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetLeaves(t *testing.T) {
	// Set up a test database

	// Create a Gin router for testing
	router := gin.Default()
	router.GET("/leaves", getLeaves)

	// Create a test request
	req, err := http.NewRequest("GET", "/leaves", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Serve the request and record the response
	router.ServeHTTP(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)
}




func TestPostLeave(t *testing.T) {
	// Set up a test database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a Gin router for testing
	router := gin.Default()
	router.POST("/leave", postLeave)

	// Create a test leave request
	newLeave := leave{
		Name:       "Jane",
		Leave_type: "casual",
		From_date:  "2023-08-23",
		To_date:    "2023-08-26",
		Team_name:  "TeamB",
		Reporter:   "Bob",
	}

	// Serialize the leave request to JSON
	jsonPayload, err := json.Marshal(newLeave)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test request with the JSON payload
	req, err := http.NewRequest("POST", "/leave", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Serve the request and record the response
	router.ServeHTTP(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, recorder.Code)

}
