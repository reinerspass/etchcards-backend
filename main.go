package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/reinerspass/waldego/internal"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

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
