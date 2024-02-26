package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func main() {
	storage, err := New()
	if err != nil {
		fmt.Print(err)
	}
	defer storage.db.Close()
	stmt, err := storage.db.Prepare(`
	INSERT INTO product (product_id, name, shelf) VALUES 
	(1, "ноутбук", "A"),
	(2, "телевизор", "A"),
	(3, "телефон", "Б"),
	(3, "телефон", "В"),
	(3, "телефон", "З"),
	(4, "системный блок", "Ж"),
	(5, "часы", "Ж"),
	(5, "часы", "А"),
	(6, "микрофон", "Ж");
	`)

	if err != nil {
		fmt.Print(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err)
	}
	stmt, err = storage.db.Prepare(`
	INSERT INTO shelf (id_shelf, product_id, base_shelf) VALUES
	("A", 1, 1),
	("A", 2, 1),
	("A", 6, 0),
	("Б", 2, 1),
	("В", 3, 0),
	("Ж", 4, 1),
	("Ж", 5, 0),
	("З", 3, 0);
	`)
	if err != nil {
		fmt.Print(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err)
	}

	stmt, err = storage.db.Prepare(`
	INSERT INTO orders (order_id, product_id, quantity) VALUES
	(10, 1, 2),
	(10, 3, 1),
	(10, 6, 1),
	(11, 2, 3),
	(14, 1, 3),
	(14, 4, 4),
	(14, 5, 1);
	`)
	if err != nil {
		fmt.Print(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err)
	}
	stmt, err = storage.db.Prepare(`SELECT * FROM product`)
	if err != nil {
		fmt.Print("01")
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Print("01")
	}
	for rows.Next() {
		var id string
		var name string
		var shelf string
		_ = rows.Scan(&id, &name, &shelf)

		fmt.Printf("%s %s %s\n", id, name, shelf)
	}
}

func New() (*Storage, error) {

	const op = "main.NewStorage"

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`DROP TABLE IF EXISTS shelf;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`DROP TABLE IF EXISTS product;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = db.Prepare(`DROP TABLE IF EXISTS orders;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
CREATE TABLE IF NOT EXISTS product(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	product_id INTEGER NOT NULL, 
	name TEXT NOT NULL,
	shelf TEXT NOT NULL);	
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS orders(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		FOREIGN KEY (product_id) REFERENCES product(product_id));
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS shelf (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		id_shelf TEXT,
		product_id INTEGER NOT NULL,
		base_shelf INTEGER CHECK(base_shelf == 0 OR base_shelf == 1),
		FOREIGN KEY (product_id) REFERENCES product(product_id));
	`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}
