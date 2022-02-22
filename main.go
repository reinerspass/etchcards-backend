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
	router.GET("/deck", deck)
	router.GET("/database", connect)
	router.Run()
	// connect()
}

func deck(c *gin.Context) {
	c.JSON(http.StatusOK, GetDeck())
}

func connect(c *gin.Context) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	response := GetData(db)
	fmt.Println("Connected!")

	response += "connected..."

	c.JSON(http.StatusOK, response)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetData(db *sql.DB) string {
	fmt.Println("get data")

	acc_string := ""

	rows, err := db.Query(`SELECT "first_name", "last_name" FROM "users"`)
	CheckError(err)

	defer rows.Close()

	for rows.Next() {
		var first_name string
		var last_name string

		err = rows.Scan(&first_name, &last_name)
		CheckError(err)

		fmt.Println(first_name, last_name)
		acc_string += first_name
		acc_string += last_name
	}

	CheckError(err)

	return acc_string
}

func GetDeck() Deck {
	return Deck{
		"Das kommt mir spanisch vor",
		[]Card{
			{"fron", "back"},
		},
	}
}

type Card struct {
	Front string
	Back  string
}

type Deck struct {
	Name  string
	Cards []Card
}
