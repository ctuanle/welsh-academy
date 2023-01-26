package main

import (
	"database/sql"
	"errors"
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

	recipes, err := app.models.Recipes.GetAll(includeMap, excludeMap)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, r, http.StatusOK, envelope{"recipes": recipes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createRecipe add an ingredient to db
// only expert can access this
func (app *application) createRecipe(w http.ResponseWriter, r *http.Request) {
	type ingredientInput struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	}

	var input struct {
		Name        string                  `json:"name"`
		CreatorId   int                     `json:"creator_id"`
		Description string                  `json:"description"`
		Ingredients map[int]ingredientInput `json:"ingredients"`
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
	v.Check(len(input.Description) > 0, "description", "recipe description can not be empty")
	v.Check(len(input.Description) < 2000, "description", "recipe description can not longer than 2000 characters")
	v.Check(input.CreatorId > 0, "creator_id", "creator_id id must be a positive integer")
	v.Check(len(input.Ingredients) > 0, "ingredients", "there must be at least one ingredient")
	v.Check(len(input.Ingredients) < 30, "ingredients", "there are way to much ingredients")

	ingredients := map[int]models.RecipeIngredient{}

	for id, info := range input.Ingredients {
		should_break := false

		v.Check(id > 0, "ingredients -> id", "ingredient id musts be a positive integer")

		// check if ingredient exists
		ing, err := app.models.Ingredients.GetById(id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				v.AddError("ingredients -> id", fmt.Sprintf("ingredient id %d does not exist", id))
				should_break = true
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		if should_break {
			// break early
			break
		}

		v.Check(info.Amount > 0, "ingredients", "ingredient amount must be positive")
		v.Check(len(info.Unit) > 0, "ingredients", "ingredient unit can not be empty")

		ingredients[id] = models.RecipeIngredient{Name: ing.Name, Amount: info.Amount, Unit: info.Unit}
	}

	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	// insert new recipes
	newRecipe := models.Recipe{Name: input.Name, Description: input.Description, CreatorId: input.CreatorId, Ingredients: ingredients}
	err = app.models.Recipes.Insert(&newRecipe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// response newly created ingredient to client
	err = app.writeJson(w, r, http.StatusCreated, envelope{"newRecipe": newRecipe}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
