package internal

import (
	"database/sql"
	"log"
	"os"
)

var database *sql.DB

func Connect() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unhandled: Error opening database: %q", err)
	}
	database = db
}

func LoadDecks() *Decks {
	rows, err := database.Query(
		`SELECT id, name 
		FROM deck`)
	if err != nil {
		log.Fatal("Unhandled: Unable to read decks: ", err)
	}
	defer rows.Close()

	var decks Decks
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(&deck.Id, &deck.Name); err != nil {
			log.Fatalf("Unhandled: error %q", err)
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
	return deck
}

func loadDeck(deckId int) *Deck {
	row := database.QueryRow(
		`SELECT id, name 
		FROM deck 
		WHERE deck.id=$1;`,
		deckId)
	var deck Deck
	if err := row.Scan(&deck.Id, &deck.Name); err != nil {
		return nil
	}
	return &deck
}

func cardIdsForDeck(deckId int) []int {
	var cardIds []int

	rows, err := database.Query(
		`SELECT card.id
		 FROM card
		 WHERE card.deck_id=$1`,
		deckId)
	if err != nil {
		log.Fatal("Unhandled: unable to load card ids: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Unhandled: Unable to scan ID for Card %q", err)
		}
		cardIds = append(cardIds, id)
	}

	return cardIds
}

func WriteDeck(deck Deck) int {
	row := database.QueryRow(
		`INSERT 
		INTO deck(name)
	 	VALUES($1)
		RETURNING id;`,
		deck.Name)
	var deckId int
	if err := row.Scan(&deckId); err != nil {
		log.Fatalf("Unhandled: Error writing Deck %q", err)
	}

	return deckId
}

func DeleteDeck(deckId int) {
	for _, cardId := range cardIdsForDeck(deckId) {
		DeleteCard(deckId, cardId)
	}
	deleteDeck(deckId)
}

func deleteDeck(deckId int) {
	row := database.QueryRow(
		`DELETE 
		FROM deck 
		WHERE id=$1
		RETURNING id;`,
		deckId)
	if err := row.Scan(&deckId); err != nil {
		log.Fatalf("Unhandled: Error deleting Deck %q", err)
	}
}

func loadCard(cardId int) Card {
	var card Card
	card.Id = cardId
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
		log.Fatal("Unhandled: Unable to load layer: ", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Type, &item.Content); err != nil {
			log.Fatalf("Unhandled: error %q", err)
		}
		items = append(items, item)
	}
	layer.Items = items

	return layer
}

func WriteCard(deckId int, card Card) int {
	frontLayerId := writeLayer()
	backLayerId := writeLayer()
	cardId := writeCard(deckId, frontLayerId, backLayerId)
	for _, item := range card.Front.Items {
		writeItem(item.Type, item.Content, frontLayerId)
	}
	for _, item := range card.Back.Items {
		writeItem(item.Type, item.Content, backLayerId)
	}
	return cardId
}

func writeCard(deckId, frontLayerId, backLayerId int) int {
	var cardId int
	row := database.QueryRow(
		`INSERT INTO card(deck_id, front_layer_id, back_layer_id)
		VALUES($1, $2, $3)
		RETURNING id;`,
		deckId, frontLayerId, backLayerId)
	if err := row.Scan(&cardId); err != nil {
		log.Fatalf("Unhandled: Error writing Card %q", err)
	}
	return cardId
}

func writeLayer() int {
	var layerId int
	row := database.QueryRow(
		`INSERT INTO layer 
		VALUES(DEFAULT)
		RETURNING id;`)
	if err := row.Scan(&layerId); err != nil {
		log.Fatalf("Unhandled: Error writing Layer %q", err)
	}
	return layerId
}

func writeItem(itemType string, content string, layerId int) int {
	var itemId int
	row := database.QueryRow(
		`INSERT INTO item(type, content, layer_id)
		VALUES($1, $2, $3)
		RETURNING id;`,
		itemType, content, layerId)
	if err := row.Scan(&itemId); err != nil {
		log.Fatalf("Unhandled: Error writing Item %q", err)
	}
	return itemId
}

func DeleteCard(deckId, cardId int) {
	if !cardInDeck(deckId, cardId) {
		log.Fatalf("Unhandled: card not contained in deck")
		return
	}
	row := database.QueryRow(
		`DELETE 
		FROM card
		CASCADE 
		WHERE id=$1
		RETURNING id;`,
		cardId)
	if err := row.Scan(&cardId); err != nil {
		log.Fatalf("Unhandled: Error writing Deck %q", err)
	}
}

func cardInDeck(deckId, cardId int) bool {
	deck := LoadDeck(deckId)
	cardInDeck := false
	for _, card := range deck.Cards {
		if cardId == card.Id {
			cardInDeck = true
			break
		}
	}
	return cardInDeck
}
