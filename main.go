package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func init_sql() string {
	return `
	-- RESET DATABASE
	DROP TABLE item;
	DROP TYPE item_type;
	DROP TABLE layer;
	DROP TABLE card;
	DROP TABLE deck;
	
	-- CREATE TABLE DECK
	CREATE TABLE deck (
		id serial PRIMARY KEY,
		name text NOT NULL
	);
	
	-- CREATE TABLE LAYER
	CREATE TABLE layer (
		id serial PRIMARY KEY
	);
	
	-- CREATE TABLE CARD
	CREATE TABLE card (
		id serial PRIMARY KEY,
		deck_id INT NOT NULL,
		CONSTRAINT fk_deck FOREIGN KEY(deck_id) REFERENCES deck(id),
		
		front_layer_id INT NOT NULL,
		CONSTRAINT fk_front_layer FOREIGN KEY(front_layer_id) REFERENCES layer(id),
		back_layer_id INT NOT NULL,
		CONSTRAINT fk_back_layer FOREIGN KEY(front_layer_id) REFERENCES layer(id)
	);
	
	-- CREATE TABLE ITEM
	CREATE TYPE item_type AS ENUM ('title', 'description', 'image');
	CREATE TABLE item (
		id serial PRIMARY KEY,
		type item_type NOT NULL,
		content text NOT NULL,
		layer_id INT NOT NULL,
		CONSTRAINT fk_layer FOREIGN KEY(layer_id) REFERENCES layer(id)
	);
	
	-- INSERT DECK
	INSERT INTO deck(name)
	VALUES('Des kommt mir spanisch vor')
	RETURNING id;
	
	-- INSERT LAYER
	INSERT INTO layer 
	VALUES(DEFAULT)
	RETURNING id;
	
	-- INSERT CARD
	INSERT INTO card(deck_id, front_layer_id, back_layer_id)
	VALUES(1, 1, 2)
	RETURNING id;
	
	-- INSERT ITEM
	INSERT INTO item(type, content, layer_id)
	VALUES('title', 'hablas', 2)
	`
}

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

func dbInit(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec(init_sql()); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}
	}
}

func readDecks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT tick FROM deck")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", name))
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "🏝 waldego")
	})
	router.GET("/db", dbFunc(db))
	router.GET("/db-init", dbInit(db))
	router.GET("/decks", readDecks(db))
	router.GET("/deck", deckMock)

	router.Run(":" + port)
}
