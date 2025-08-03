package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Cell holds the data for a single grid cell
type Cell struct {
	Value   string
	Classes string // To hold CSS classes like "thick-bottom"
}

// PageData will hold the data for our template
type PageData struct {
	GameTitle string
	Rows      [][]Cell // The grid is now a slice of Cell slices
}

// gridHandler parses URL parameters, creates the grid, and executes the template
func gridHandler(w http.ResponseWriter, r *http.Request) {
	boardMultiplier := 3
	// Get 'size' from query params, default to 3 (for a 9x9 grid)
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 3 // Default size
	}

	gridSize := size * boardMultiplier

	// Create the grid data (a slice of Cell slices)
	grid := make([][]Cell, gridSize)
	for i := range grid {
		grid[i] = make([]Cell, gridSize)
		for j := range grid[i] {
			var classes []string

			// Check if it's a third row (and not the last row)
			if (i+1)%boardMultiplier == 0 && i+1 < gridSize {
				classes = append(classes, "thick-bottom")
			}
			// Check if it's a third column (and not the last column)
			if (j+1)%boardMultiplier == 0 && j+1 < gridSize {
				classes = append(classes, "thick-right")
			}

			grid[i][j] = Cell{
				Value:   fmt.Sprintf("%d,%d", i+1, j+1),
				Classes: strings.Join(classes, " "), // e.g., "thick-bottom thick-right"
			}
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
		GameTitle: fmt.Sprintf("%dx%d Grid", gridSize, gridSize),
		Rows:      grid,
	}

	// Execute the template, passing it the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/grid", gridHandler)
	log.Println("Server starting on :8080... Access http://localhost:8080/grid?size=3")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
