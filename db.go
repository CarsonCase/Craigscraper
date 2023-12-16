package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
)

type Listing struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func SetupDB() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./db/listings.db")
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS listings (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT,
				price TEXT,
				link TEXT
			)
		`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertListing(db *sql.DB, listing Listing) error {
	stmt, err := db.Prepare("INSERT INTO listings (title, price, link) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(listing.Title, listing.Price, listing.Link)

	if err != nil {
		return err
	}
	return nil
}

func getAllListings(db *sql.DB) ([]Listing, error) {
	rows, err := db.Query("SELECT * FROM listings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Price, &listing.Link)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}
	return listings, nil
}
