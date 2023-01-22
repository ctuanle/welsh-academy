package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"ctuanle.ovh/welsh-academy/internal/models"
)

const version = "1.0.0"

type config struct {
	port int
}

type application struct {
	config      config
	logger      *log.Logger
	ingredients *models.IngredientModel
	recipes     *models.RecipeModel
	favorites   *models.FavoriteModel
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:      cfg,
		logger:      logger,
		ingredients: &models.IngredientModel{Ingredients: models.Ingredients},
		recipes:     &models.RecipeModel{Recipes: models.Recipes},
		favorites:   &models.FavoriteModel{Favorites: models.Favorites},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	logger.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
