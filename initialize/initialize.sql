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