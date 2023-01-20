package main

import (
	"fmt"
	"net/http"
)

// listRecipes list all existing recipes
// both expert and user can access this
func (app *application) listRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list all recipes here")
}

// createRecipe add an ingredient to db
// only expert can access this
func (app *application) createRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "add an recipe, expert only")
}
