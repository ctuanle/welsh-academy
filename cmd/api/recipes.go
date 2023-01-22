package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ctuanle.ovh/welsh-academy/internal/models"
	"ctuanle.ovh/welsh-academy/internal/validator"
)

// listRecipes list all existing recipes
// both expert and user can access this
// support including/excluding ingredient(s)
func (app *application) listRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	include := strings.Split(query.Get("include"), ",")
	exclude := strings.Split(query.Get("exclude"), ",")

	includeMap := make(map[int]struct{})
	excludeMap := make(map[int]struct{})

	for _, rid := range include {
		if rid != "" {
			_rid, err := strconv.Atoi(rid)
			if err != nil || _rid < 1 {
				app.errorResponse(w, r, http.StatusBadRequest, "Invalid Ingredient ID (include)")
				return
			}
			includeMap[_rid] = struct{}{}
		}
	}

	for _, rid := range exclude {
		if rid != "" {
			_rid, err := strconv.Atoi(rid)
			if err != nil || _rid < 1 {
				app.errorResponse(w, r, http.StatusBadRequest, "Invalid Ingredient ID (exclude)")
				return
			}

			if _, ok := includeMap[_rid]; ok {
				app.errorResponse(w, r, http.StatusBadRequest, fmt.Sprintf("Cannot include and exclude ingredient with ID %d", _rid))
				return
			}

			excludeMap[_rid] = struct{}{}
		}

	}

	recipes, _ := app.recipes.GetAll(includeMap, excludeMap)
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
	err := app.readBodyToJSON(w, r, &input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	// validate input
	v := validator.New()
	v.Check(len(input.Name) > 0, "name", "recipe name can not be empty")
	v.Check(len(input.Name) < 100, "name", "recipe name can not longer than 100 characters")
	v.Check(len(input.Description) > 0, "name", "recipe description can not be empty")
	v.Check(len(input.Description) > 2000, "name", "recipe description can not longer than 2000 characters")
	v.Check(input.Creator > 0, "creator", "creator id must be a positive integer")
	for _, ing := range input.Ingredients {
		v.Check(ing.ID > 0, "ingredients", "ingredient id must be a positive integer")
		v.Check(ing.Amount > 0, "ingredients", "ingredient amount must be positive")
		v.Check(len(ing.Unit) > 0, "ingredients", "ingredient unit can not be empty")
	}

	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
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
