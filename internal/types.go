package internal

type Deck struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type Decks struct {
	Decks []Deck `json:"decks"`
}

type Cards struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	Front Layer `json:"front"`
	Back  Layer `json:"back"`
}

type Layer struct {
	Items []Item `json:"items"`
}

type Item struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
