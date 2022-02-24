package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	router.POST("/decks", func(c *gin.Context) {
		var deck internal.Deck
		if err := c.BindJSON(&deck); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("received valid object? ", deck)
		var deckId = internal.WriteDeck(deck)

		c.JSON(http.StatusOK, gin.H{"status": "success", "deck_id": deckId})
	})
	router.GET("/decks/:deck_id", func(c *gin.Context) {
		var deck *internal.Deck
		deck_id, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
			// log.Fatal("unable to parse parameter: ", err)
		}
		deck = internal.LoadDeck(deck_id)
		if deck == nil {
			c.JSON(http.StatusNotFound, "404 Deck Not Found")
		} else {
			c.JSON(http.StatusOK, deck)
		}
	})

	router.Run(":" + port)
}
