package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"ctuanle.ovh/welsh-academy/internal/models"
	"ctuanle.ovh/welsh-academy/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// listFavorites list all favorite recipes of a user
func (app *application) listFavorites(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get user id
	uid, err := strconv.Atoi(params.ByName("uid"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	favorites, err := app.models.Favorites.GetAll(uid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, r, http.StatusOK, envelope{"favorites": favorites}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// flagFavoriteRecipe flags recipe as favorite
func (app *application) flagFavoriteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get user id
	uid, err := strconv.Atoi(params.ByName("uid"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		RecipeId int `json:"recipe_id"`
	}

	err = app.readBodyToJSON(w, r, &input)
	if err != nil {
		// bad request
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// validate input
	v := validator.New()
	v.Check(input.RecipeId > 0, "recipe_id", "recipe_id must be a positive integer")
	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	newFav := models.Favorite{UserId: uid, RecipeId: input.RecipeId}
	err = app.models.Favorites.Insert(&newFav)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, r, http.StatusCreated, envelope{"newFavorite": newFav}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// unflagFavoriteRecipe unflags a favorite recipe
func (app *application) unflagFavoriteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get recipe id
	fid, err := strconv.Atoi(params.ByName("fid"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Favorites.Remove(fid)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJson(w, r, http.StatusOK, envelope{"message": "Deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
