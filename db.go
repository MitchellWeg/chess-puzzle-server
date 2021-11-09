package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/MonetDB/MonetDB-Go/src"
	"github.com/joho/godotenv"
)

type Database struct {
	handle *sql.DB
}

type Puzzle struct {
	Puzzleid        string
	Fen             string
	Moves           string
	Rating          string
	Popularity      string
	Ratingdeviation string
	Nbplays         string
	Themes          string
	Gameurl         string
}

func (d *Database) init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("username")
	password := os.Getenv("password")
	host := os.Getenv("host")
	port := os.Getenv("port")
	database := os.Getenv("database")

	db, err := sql.Open("monetdb", fmt.Sprintf("%s:%s@%s:%s/%s", username, password, host, port, database))

	if err != nil {
		panic(err)
	}

	d.handle = db
}

func (d *Database) query(query string) ([]Puzzle, error) {
	rows, err := d.handle.Query(query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var puzzles []Puzzle

	for rows.Next() {
		var _puzzle Puzzle
		if err := rows.Scan(&_puzzle.Puzzleid, &_puzzle.Fen, &_puzzle.Moves,
			&_puzzle.Rating, &_puzzle.Popularity, &_puzzle.Ratingdeviation, &_puzzle.Nbplays, &_puzzle.Themes, &_puzzle.Gameurl); err != nil {
			return puzzles, err
		}

		puzzles = append(puzzles, _puzzle)
	}
	if err = rows.Err(); err != nil {
		return puzzles, err
	}

	return puzzles, nil
}
