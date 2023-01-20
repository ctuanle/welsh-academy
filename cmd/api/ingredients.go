package main

import (
	"fmt"
	"net/http"
)

// listIngredients list all existing ingredients
// both expert and user can access this
func (app *application) listIngredients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list all ingredients here")
}

// createIngredient add an ingredient to db
// only expert can access this
func (app *application) createIngredient(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "add an ingredient, expert only")
}
