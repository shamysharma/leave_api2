package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getId(c *gin.Context) {
	id := c.Param("id")

	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, leave_type, from_date, to_date, team_name, sick_leaves_file, reporter FROM leaves WHERE id= $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var leaves []leave_for_get
	for rows.Next() {
		var a leave_for_get
		err := rows.Scan(&a.Id, &a.Name, &a.Leave_type, &a.From_date, &a.To_date, &a.Team_name, &a.Sick_leaves_file, &a.Reporter)
		if err != nil {
			log.Fatal(err)
		}
		leaves = append(leaves, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, leaves)

}

func getLeaves(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, leave_type, from_date, to_date, team_name, sick_leaves_file, reporter FROM leaves")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var leaves []leave_for_get
	for rows.Next() {
		var a leave_for_get
		err := rows.Scan(&a.Id, &a.Name, &a.Leave_type, &a.From_date, &a.To_date, &a.Team_name, &a.Sick_leaves_file, &a.Reporter)
		if err != nil {
			log.Fatal(err)
		}
		leaves = append(leaves, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, leaves)
}

func postLeave(c *gin.Context) {

	var asLeave leave
	if err := c.ShouldBind(&asLeave); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if asLeave.Leave_type == "Sick Leave" {
		// Check if the file was uploaded
		file, err := c.FormFile("file")
		if file == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload is required for Sick Leave"})
			return
		}
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"file": "file successfully uploaded"})
		}

		tempfilePath := "./documents/temp/" + file.Filename

		if err := c.SaveUploadedFile(file, tempfilePath); err != nil {

			c.String(http.StatusInternalServerError, "Failed to save file")

			return

		}

		_, err = db.Exec("INSERT INTO leaves (name, leave_type, from_date, to_date, team_name, sick_leaves_file, reporter) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			asLeave.Name, asLeave.Leave_type, asLeave.From_date, asLeave.To_date, asLeave.Team_name, file.Filename, asLeave.Reporter)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, asLeave)
	} else {
		_, err1 := db.Exec("INSERT INTO leaves (name, leave_type, from_date, to_date, team_name, sick_leaves_file, reporter) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			asLeave.Name, asLeave.Leave_type, asLeave.From_date, asLeave.To_date, asLeave.Team_name, "null", asLeave.Reporter)
		if err1 != nil {
			log.Fatal(err1)
		}
		c.JSON(http.StatusCreated, asLeave)
	}

}

func getFile(c *gin.Context) {
	id := c.Param("id")

	var fileName string
	err := db.QueryRow("SELECT sick_leaves_file FROM leaves WHERE id = $1", id).Scan(&fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file information"})
		return
	}

	if fileName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	file, err := os.Open(fmt.Sprintf("./documents/temp/%s", fileName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file information"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
}
