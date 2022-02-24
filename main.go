package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/reinerspass/waldego/internal"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func deckMock(c *gin.Context) {
	c.String(http.StatusOK, `{
		"decks": [
			{
				"id": 1,
				"name": "Des hört sich spanisch ah"
			},
			{
				"id": 2,
				"name": "Des hört sich französich ah"
			},
			{
				"id": 3,
				"name": "Des hört sich dütsch ah"
			},
			{
				"id": 4,
				"name": "Des hört sich englisch ah"
			}
		]
	}`)
}

func dbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	internal.Connect()

	router := gin.New()
	logger := gin.Logger()
	router.Use(logger)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ha")
	})
	router.GET("/test", func(c *gin.Context) {
		internal.Initialize()
		c.JSON(http.StatusOK, "ha")
	})
	router.GET("/decks", func(c *gin.Context) {
		var decks = internal.LoadDecks()
		c.JSON(http.StatusOK, decks)
	})
	router.GET("/deck", deckMock)
	router.GET("/cards", func(c *gin.Context) {
		var cards = internal.LoadsCards(1)
		fmt.Printf("wasch da los", cards)
		c.JSON(http.StatusOK, cards)
	})

	router.Run(":" + port)
}
