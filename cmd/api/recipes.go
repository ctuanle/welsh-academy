package main

import (
	"encoding/json"
	"net/http"

	"ctuanle.ovh/welsh-academy/internal/models"
)

// listRecipes list all existing recipes
// both expert and user can access this
func (app *application) listRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, _ := app.recipes.GetAll()
	err := app.writeJson(w, r, http.StatusOK, envelope{"recipes": recipes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createRecipe add an ingredient to db
// only expert can access this
func (app *application) createRecipe(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string                    `json:"name"`
		Creator     int                       `json:"creator"`
		Description string                    `json:"description"`
		Ingredients []models.RecipeIngredient `json:"ingredients"`
	}

	// decode body content into input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// insert new recipes
	newRecipe, _ := app.recipes.Insert(input.Name, input.Description, input.Creator, input.Ingredients)

	// response newly created ingredient to client
	err = app.writeJson(w, r, http.StatusCreated, envelope{"newRecipe": newRecipe}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
