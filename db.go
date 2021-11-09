package main

import (
	"database/sql"

	_ "github.com/MonetDB/MonetDB-Go/src"
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
	db, err := sql.Open("monetdb", "monetdb:monetdb@localhost:50000/demo")

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
