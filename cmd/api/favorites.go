package main

import (
	"fmt"
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

	fmt.Fprintf(w, "list favorite recipes of user %d\n", uid)
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

	fmt.Fprintf(w, "flag recipe as favorite for user %d\n", uid)
}

// unflagFavoriteRecipe unflags a favorite recipe
func (app *application) unflagFavoriteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get user id
	uid, err := strconv.ParseInt(params.ByName("uid"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// get recipe id
	rid, err := strconv.ParseInt(params.ByName("rid"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	fmt.Fprintf(w, "unflag recipe %d for user %d\n", rid, uid)
}
