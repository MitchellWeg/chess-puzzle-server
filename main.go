package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	db Database
}

func main() {
	api := API{}

	api.serve()
}

func (a *API) serve() {
	a.db = Database{}
	a.db.init()

	r := mux.NewRouter()
	r.HandleFunc("/puzzle", a.puzzleHandler)
	r.HandleFunc("/puzzle/rating", a.puzzleByRatingHandler)
	r.HandleFunc("/puzzles", a.multiplePuzzlesHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

/*
* Route: <URL>:<PORT>/puzzle?minRating={min}&maxRating={max}&theme={theme}
* Gets a random puzzle.
* All parameters are optional and can be omitted any way you please.
 */
func (a *API) puzzleHandler(w http.ResponseWriter, r *http.Request) {
	min := 0
	max := 3000
	theme := ""

	vals := r.URL.Query()
	_min, minOk := vals["minRating"]
	_max, maxOk := vals["maxRating"]
	_theme, themeOk := vals["theme"]

	if minOk {
		a, err := strconv.Atoi(_min[0])

		if err != nil {
			w.Write([]byte("Query not complete: min is not an int"))
			return
		}

		min = a
	}

	if maxOk {
		a, err := strconv.Atoi(_max[0])

		if err != nil {
			w.Write([]byte("Query not complete: max is not an int"))
			return
		}

		max = a
	}

	if themeOk {
		parsedTheme, err := filterTheme(strings.ToLower(_theme[0]))

		if err != nil {
			w.Write([]byte(fmt.Sprintf("Theme '%s' is not a valid theme", _theme[0])))
			return
		}

		theme = parsedTheme
	}

	q := ""

	if theme != "" {
		q = fmt.Sprintf("SELECT * FROM puzzles WHERE LOWER(themes) LIKE '%%%s%%' AND rating BETWEEN %d AND %d SAMPLE 1", theme, min, max)
	} else {
		q = fmt.Sprintf("SELECT * FROM puzzles WHERE rating BETWEEN %d AND %d SAMPLE 1", min, max)
	}

	puzzles, err := a.db.query(q)

	if err != nil {
		panic(err)
	}

	j := a.serializeToJson(puzzles)

	w.Write(j)
}

/*
* Route: <URL>:<PORT>/puzzle/rating?min={min}&max={max}
* Gets a random puzzle by rating between a min and a max.
 */
func (a *API) puzzleByRatingHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	_min, minOk := vals["min"]
	_max, maxOk := vals["max"]

	if !minOk {
		w.Write([]byte("Query not complete: min missing"))
		return
	}

	if !maxOk {
		w.Write([]byte("Query not complete: max missing"))
		return
	}

	min, err := strconv.Atoi(_min[0])

	if err != nil {
		w.Write([]byte("Query not complete: min is not an int"))
		return
	}

	max, err := strconv.Atoi(_max[0])

	if err != nil {
		w.Write([]byte("Query not complete: max is not an int"))
		return
	}

	puzzles, err := a.db.query(fmt.Sprintf("SELECT * FROM puzzles WHERE rating BETWEEN %d AND %d SAMPLE 1000", min, max))

	if err != nil {
		panic(err)
	}

	j := a.serializeToJson(puzzles)

	w.Write(j)
}

/*
* Route: <URL>:<PORT>/puzzles
* Gets 1000 random puzzles.
 */
func (a *API) multiplePuzzlesHandler(w http.ResponseWriter, r *http.Request) {
	puzzles, err := a.db.query("SELECT * FROM puzzles SAMPLE 1000")

	if err != nil {
		panic(err)
	}

	j := a.serializeToJson(puzzles)

	w.Write(j)
}

func filterTheme(input string) (string, error) {
	fmt.Println(input)
	validThemes := []string{"mate", "matein1", "matein2", "matein3", "matein4", "matein5", "advantage", "defensivemove", "hangingpiece", "middlegame", "verylong", "short", "trappedpiece", "endgame", "crushing", "motion", "fork", "opening", "attraction", "deflection", "long", "clearance", "kingsideattack", "sacrifice", "master", "arabianmate", "backrankmate", "exposedking", "rookendgame", "equality", "xrayattack", "pin", "discoveredattack"}

	if contains(validThemes, input) {
		return input, nil
	}

	return "", errors.New(fmt.Sprintf("Theme %s not found", input))
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (a *API) serializeToJson(p []Puzzle) []byte {
	j, err := json.Marshal(p)

	if err != nil {
		panic(err)
	}

	return j
}
