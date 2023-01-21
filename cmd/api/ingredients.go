package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ctuanle.ovh/welsh-academy/internal/data"
)

var ingredients = []data.Ingredient{
	{
		ID:        "farine",
		Name:      "Farine",
		CreatedAt: time.Now(),
	},
	{
		ID:        "fromage",
		Name:      "Fromage",
		CreatedAt: time.Now(),
	},
	{
		ID:        "piment",
		Name:      "Piment",
		CreatedAt: time.Now(),
	},
}

// listIngredients list all existing ingredients
// both expert and user can access this
func (app *application) listIngredients(w http.ResponseWriter, r *http.Request) {
	err := app.writeJson(w, r, http.StatusOK, envelope{"ingredients": ingredients}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}

// createIngredient add an ingredient to db
// only expert can access this
func (app *application) createIngredient(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	// decode body content into input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	newIngredient := data.Ingredient{
		ID:   strings.ToLower(input.Name),
		Name: input.Name,
	}

	ingredients = append(ingredients, newIngredient)

	err = app.writeJson(w, r, http.StatusCreated, nil, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
