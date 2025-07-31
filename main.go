package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// PageData will hold the data for our template
type PageData struct {
	Rows [][]string // A 2D slice to represent the grid
}

// gridHandler parses URL parameters, creates the grid, and executes the template
func gridHandler(w http.ResponseWriter, r *http.Request) {
	// Get 'n' (rows) from query params, default to 10
	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		n = 10 // Default number of rows
	}

	// Get 'm' (columns) from query params, default to 10
	mStr := r.URL.Query().Get("m")
	m, err := strconv.Atoi(mStr)
	if err != nil || m <= 0 {
		m = 10 // Default number of columns
	}

	// Create the grid data (a slice of slices)
	grid := make([][]string, n)
	for i := range grid {
		grid[i] = make([]string, m)
		for j := range grid[i] {
			// You can put any data you want in the cell.
			// Here, we just put the coordinates.
			grid[i][j] = strconv.Itoa(i+1) + "," + strconv.Itoa(j+1)
		}
	}

	// Parse our template file
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the data structure to pass to the template
	data := PageData{
		Rows: grid,
	}

	// Execute the template, passing it the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/grid", gridHandler)
	log.Println("Server starting on :8080... Access http://localhost:8080/grid")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
