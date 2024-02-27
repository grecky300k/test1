package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}
type output struct {
	product_id  string
	name        string
	id_shelf    string
	extra_shelf string
	quantity    string
	order_id    string
}

func main() {
	const op = "main"
	storage, err := New()
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	defer storage.db.Close()

	args := os.Args
	s := strings.Join(args[1:], "")
	args = strings.Split(s, ",")

	storage = Insertion(storage)

	stmt, err := storage.db.Prepare(`SELECT product.product_id, product.name, shelf.id_shelf, shelf.extra_shelf, orders.quantity, orders.order_id
	FROM product
	JOIN orders ON product.product_id = orders.product_id
	JOIN shelf ON product.product_id = shelf.product_id 
	ORDER BY shelf.id_shelf, orders.order_id, product.product_id
	  ;`)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query()

	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	outputs := []output{}

	for rows.Next() {
		out := output{}

		err = rows.Scan(&out.product_id, &out.name, &out.id_shelf, &out.extra_shelf, &out.quantity, &out.order_id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		outputs = append(outputs, out)
	}
	title := ""
	for _, out := range outputs {
		if title == "" || out.id_shelf != title {
			title = out.id_shelf
			fmt.Printf("=== Стеллаж %v\n", title)
			if out.extra_shelf != "0" {
				fmt.Printf("%v, (id=%v)\nзаказ %v, %v шт\nдоп стеллаж: %v \n\n", out.name, out.product_id, out.order_id, out.quantity, out.extra_shelf)
			} else {
				fmt.Printf("%v, (id=%v)\nзаказ %v, %v шт\n\n", out.name, out.product_id, out.order_id, out.quantity)
			}
		} else if out.id_shelf == title {
			if out.extra_shelf != "0" {
				fmt.Printf("%v, (id=%v)\nзаказ %v, %v шт\nдоп стеллаж: %v \n\n", out.name, out.product_id, out.order_id, out.quantity, out.extra_shelf)
			} else {
				fmt.Printf("%v, (id=%v)\nзаказ %v, %v шт\n\n", out.name, out.product_id, out.order_id, out.quantity)
			}
		}

	}
}

func New() (*Storage, error) {

	const op = "main.NewStorage"

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS product(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL, 
		name TEXT NOT NULL);	
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
		extra_shelf TEXT,
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

func Insertion(storage *Storage) *Storage {
	const op = "main.Insertion"
	stmt, err := storage.db.Prepare(`
	INSERT INTO product (product_id, name) VALUES 
	(1, "ноутбук"),
	(2, "телевизор"),
	(3, "телефон"),
	(4, "системный блок"),
	(5, "часы"),
	(6, "микрофон");
	`)

	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = storage.db.Prepare(`
	INSERT INTO shelf (id_shelf, product_id, extra_shelf) VALUES
	("A", 1, 0),
	("A", 2, 0),
	("Б", 3, "В,З"),
	("Ж", 4, 0),
	("Ж", 5, "A"),
	("Ж", 6, 0);
	`)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = storage.db.Prepare(`
	INSERT INTO orders (order_id, product_id, quantity) VALUES
	(10, 1, 2),
	(10, 3, 1),
	(10, 6, 1),
	(11, 2, 3),
	(14, 1, 3),
	(14, 4, 4),
	(15, 5, 1);
	`)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return storage
}
