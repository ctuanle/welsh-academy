package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// listFavorites list all favorite recipes of a user
func (app *application) listFavorites(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get user id
	uid, err := strconv.ParseInt(params.ByName("uid"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	favorites, _ := app.favorites.GetAll(int(uid))
	err = app.writeJson(w, r, http.StatusOK, envelope{"favorites": favorites}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// flagFavoriteRecipe flags recipe as favorite
func (app *application) flagFavoriteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get user id
	uid, err := strconv.ParseInt(params.ByName("uid"), 10, 64)
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

	newFav, _ := app.favorites.Insert(int(uid), input.RecipeId)
	err = app.writeJson(w, r, http.StatusCreated, envelope{"newFavorite": newFav}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// unflagFavoriteRecipe unflags a favorite recipe
func (app *application) unflagFavoriteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get recipe id
	fid, err := strconv.ParseInt(params.ByName("fid"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_ = app.favorites.Remove(int(fid))
	err = app.writeJson(w, r, http.StatusOK, envelope{"message": "Deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
