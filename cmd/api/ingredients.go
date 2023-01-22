package main

import (
	"encoding/json"
	"net/http"
)

// listIngredients list all existing ingredients
// both expert and user can access this
func (app *application) listIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, _ := app.ingredients.GetAll()
	err := app.writeJson(w, r, http.StatusOK, envelope{"ingredients": ingredients}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createIngredient add an ingredient to db
// only expert can access this
func (app *application) createIngredient(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Creator int    `json:"creator"`
	}

	// decode body content into input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// insert new ingredient
	newIngredientId, _ := app.ingredients.Insert(input.Name, input.Creator)

	// response newly created ingredient to client
	err = app.writeJson(w, r, http.StatusCreated, envelope{"newIngredient": newIngredientId}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
