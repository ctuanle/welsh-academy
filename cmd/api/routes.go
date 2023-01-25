package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// custom notfound and method-not-allowed handler
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// root route
	router.GET("/", toHandle(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Welsh Academy.")
	}))

	// ingredients routes
	router.GET("/ingredients", toHandle(app.listIngredients))
	router.POST("/ingredients", toHandle(app.createIngredient))

	//recipes routes
	router.GET("/recipes", toHandle(app.listRecipes))
	router.POST("/recipes", toHandle(app.createRecipe))

	// favorite recipes routes
	// router.GET("/users/:uid/favorites", toHandle(app.listFavorites))
	// router.POST("/users/:uid/favorites", toHandle(app.flagFavoriteRecipe))
	// router.DELETE("/users/:uid/favorites/:fid", toHandle(app.unflagFavoriteRecipe))
	router.HandlerFunc(http.MethodGet, "/users/:uid/favorites", app.listFavorites)
	router.HandlerFunc(http.MethodPost, "/users/:uid/favorites", app.flagFavoriteRecipe)
	router.HandlerFunc(http.MethodDelete, "/users/:uid/favorites/:fid", app.unflagFavoriteRecipe)

	return app.logRequest(enableCORS(router))
}

func toHandle(h func(w http.ResponseWriter, r *http.Request)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h(w, r)
	}
}
