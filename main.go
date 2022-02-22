package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "markus"
	password = ""
	dbname   = "markus"
)

func main() {
	router := gin.Default()
	router.GET("/deck/:slug", deck)
	router.Run()
	// connect()
}

func deck(c *gin.Context) {
	slug := c.Param("slug")
	c.JSON(http.StatusOK, gin.H{
		"slug": slug,
		"name": "Das kommt mir spanisch vor",
	})
}

func connect() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	GetData(db)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetData(db *sql.DB) {
	fmt.Println("get data")

	rows, err := db.Query(`SELECT "first_name", "last_name" FROM "users"`)
	CheckError(err)

	defer rows.Close()

	for rows.Next() {
		var first_name string
		var last_name string

		err = rows.Scan(&first_name, &last_name)
		CheckError(err)

		fmt.Println(first_name, last_name)
	}

	CheckError(err)
}
