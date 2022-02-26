package main

import (
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
		c.JSON(http.StatusOK, "🔪 ETCH CARDS")
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
		var deckId = internal.WriteDeck(deck)
		c.JSON(http.StatusOK, gin.H{"status": "success", "deck_id": deckId})
	})

	router.GET("/decks/:deck_id", func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		deck := internal.LoadDeck(deckId)
		if deck == nil {
			c.JSON(http.StatusNotFound, "404 Deck Not Found")
		} else {
			c.JSON(http.StatusOK, deck)
		}
	})

	router.DELETE("/decks/:deck_id", func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		internal.DeleteDeck(deckId)
		c.JSON(http.StatusOK, gin.H{"status": "success", "deck_id": deckId})
	})

	router.POST("/decks/:deck_id/cards", func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var card internal.Card
		if err := c.BindJSON(&card); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var cardId = internal.WriteCard(deckId, card)

		c.JSON(http.StatusOK, gin.H{"status": "success", "card_id": cardId})
	})

	router.DELETE("/decks/:deck_id/cards/:card_id", func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cardId, err := strconv.Atoi(c.Param("card_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		internal.DeleteCard(deckId, cardId)
		c.JSON(http.StatusOK, gin.H{"status": "success", "card_id": cardId})
	})

	router.Run(":" + port)
}
