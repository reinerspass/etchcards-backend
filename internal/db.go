package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var database *sql.DB

func Connect() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: ", err)
	}
	database = db
}

func Initialize() {
	if _, err := database.Exec(init_sql()); err != nil {
		log.Fatal("unable to initialize database: ", err)
		return
	}
}

func LoadDecks() *Decks {
	rows, err := database.Query("SELECT id, name FROM deck")
	if err != nil {
		log.Fatal("unable to read decks: ", err)
	}
	defer rows.Close()

	var decks Decks
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(&deck.Id, &deck.Name); err != nil {
			log.Fatalf("error %q", err)
		}
		decks.Decks = append(decks.Decks, deck)
	}

	return &decks
}

func LoadDeck(deckId int) *Deck {
	// return nil
	deck := loadDeck(deckId)
	for _, cardId := range cardIdsForDeck(deckId) {
		deck.Cards = append(deck.Cards, loadCard(cardId))
	}
	if deck == nil {
		return nil
	}
	if len(deck.Cards) == 0 {
		return nil
	}
	return deck
}

func WriteDeck(deck Deck) int {
	var deckId int

	row := database.QueryRow(
		`INSERT INTO deck(name)
	 	 VALUES($1)
		 RETURNING id;`, deck.Name)

	if err := row.Scan(&deckId); err != nil {
		log.Fatalf("error error error %q", err)
	}

	return deckId
}

func loadDeck(deckId int) *Deck {
	var deck Deck
	row := database.QueryRow(`SELECT id, name FROM deck where deck.id=$1;`, deckId)

	if err := row.Scan(&deck.Id, &deck.Name); err != nil {
		log.Fatalf("error error error %q", err)
	}

	return &deck
}

func loadCard(cardId int) Card {
	var card Card
	card.Front = loadLayer(cardId, true)
	card.Back = loadLayer(cardId, false)
	return card
}

func loadLayer(cardId int, front bool) Layer {
	var layer Layer
	var sql string

	if front {
		sql = `SELECT item.type, item.content
		FROM item
		INNER JOIN layer
		ON item.layer_id=layer.id
		INNER JOIN card
		on layer.id=card.front_layer_id
		where card.id=$1;`
	} else {
		sql = `SELECT item.type, item.content
		FROM item
		INNER JOIN layer
		ON item.layer_id=layer.id
		INNER JOIN card
		on layer.id=card.back_layer_id
		where card.id=$1;`
	}

	rows, err := database.Query(
		sql, cardId)
	if err != nil {
		log.Fatal("unable to load layer: ", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item

		if err := rows.Scan(&item.Type, &item.Content); err != nil {
			log.Fatalf("error %q", err)
		}

		fmt.Println("loading on the line", item)

		items = append(items, item)
	}
	layer.Items = items

	return layer
}

func cardIdsForDeck(deckId int) []int {
	var cardIds []int

	rows, err := database.Query(
		`SELECT card.id
		 FROM card
		 WHERE card.deck_id=$1`,
		deckId)
	if err != nil {
		log.Fatal("unable to cards: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("error %q", err)
		}

		cardIds = append(cardIds, id)
	}

	return cardIds
}

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
	VALUES('Des kommt mir französisch vor')
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
	VALUES('title', 'hablas', 2);
	
	INSERT INTO item(type, content, layer_id)
	VALUES('title', 'reden', 1);
	
	
	-- GET DECKS
	SELECT * FROM deck;
	
	-- GET ITEMS FOR DECK
	SELECT * 
	FROM card
	INNER JOIN layer 
	ON layer.id=card.front_layer_id
	OR layer.id=card.back_layer_id
	INNER JOIN item
	ON layer.id=item.layer_id
	WHERE card.deck_id=1;
	
	-- GET CARDS FOR DECK
	SELECT card.id
	FROM card
	WHERE card.deck_id=1;
	
	--GET FRONT ITEMS OF CARD
	SELECT item.type, item.content
	FROM item
	INNER JOIN layer
	ON item.layer_id=layer.id
	INNER JOIN card
	on layer.id=card.front_layer_id
	where card.id=1;
	
	--GET FRONT ITEMS OF CARD
	SELECT item.type, item.content
	FROM item
	INNER JOIN layer
	ON item.layer_id=layer.id
	INNER JOIN card
	on layer.id=card.back_layer_id
	where card.id=1;
	`
}
