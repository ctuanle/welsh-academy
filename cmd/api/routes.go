package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// root route
	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Welsh Academy.")
	})

	// ingredients routes
	router.HandlerFunc(http.MethodGet, "/ingredients", app.listIngredients)
	router.HandlerFunc(http.MethodPost, "/ingredients", app.createIngredient)

	//recipes routes
	router.HandlerFunc(http.MethodGet, "/recipes", app.listRecipes)
	router.HandlerFunc(http.MethodPost, "/recipes", app.createRecipe)

	// favorite recipes routes
	router.HandlerFunc(http.MethodGet, "/users/:uid/favorites", app.listFavorites)
	router.HandlerFunc(http.MethodPost, "/users/:uid/favorites", app.flagFavoriteRecipe)
	router.HandlerFunc(http.MethodDelete, "/users/:uid/favorites/:rid", app.unflagFavoriteRecipe)

	return router
}
