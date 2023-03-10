package main

import (
	"net/http"
	"strings"

	"ctuanle.ovh/welsh-academy/internal/models"
	"ctuanle.ovh/welsh-academy/internal/validator"
	"github.com/lib/pq"
)

// listIngredients list all existing ingredients
// both expert and user can access this
func (app *application) listIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := app.models.Ingredients.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, r, http.StatusOK, envelope{"ingredients": ingredients}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createIngredient add an ingredient to db
// only expert can access this
func (app *application) createIngredient(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		CreatorId int    `json:"creator_id"`
	}

	// decode body content into input
	err := app.readBodyToJSON(w, r, &input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.Name = strings.TrimSpace(input.Name)

	// validate input
	v := validator.New()
	v.Check(len(input.Name) > 0, "name", "must not be empty")
	v.Check(input.CreatorId > 0, "creator_id", "creator_id id must be a positive integer")
	v.Check(len(input.Name) < 100, "name", "must be less than 100 characters")

	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	newIng := models.Ingredient{
		Name:      input.Name,
		CreatorId: input.CreatorId,
	}

	err = app.models.Ingredients.Insert(&newIng)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "foreign_key_violation" {
				app.errorResponse(w, r, http.StatusBadRequest, "creator id does not exist")
				return
			}
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	/* headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/ingredients/%d", newIng.ID)) */

	// response newly created ingredient to client
	err = app.writeJson(w, r, http.StatusCreated, envelope{"ingredient": newIng}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
