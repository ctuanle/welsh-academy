package main

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"ctuanle.ovh/welsh-academy/internal/models"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func newTestApplication() *application {
	return &application{
		logger: log.New(io.Discard, "", 0),
		models: models.NewMockModels(),
	}
}

func newTestServer() *http.Server {
	return &http.Server{
		Handler: newTestApplication().routes(),
	}
}

func newTestAPI(t *testing.T) *tdhttp.TestAPI {
	return tdhttp.NewTestAPI(t, newTestServer().Handler)
}

func TestListIngredientsHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	testAPI.Get("/ingredients").
		Name("List all ingredients").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Ingredient{
			"ingredients": models.MockedIngredients,
		})
}

func TestCreateIngredientHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	body := map[string]interface{}{
		"name":       "Test Name",
		"creator_id": 1,
	}

	expectedNewID := len(models.MockedIngredients) + 1

	testAPI.PostJSON("/ingredients", body).
		Name("Create Ingredient").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(td.JSON(`
			{
				"ingredient": {
					"id": $id,
					"name": "Test Name",
					"creator_id": 1,
					"created": "$created"
				}
			}
		`,
			td.Tag("id", expectedNewID),
			td.Tag("created", td.Between(testAPI.SentAt(), time.Now())),
		))

	invalidBody := map[string]interface{}{
		"name":       "",
		"creator_id": 0,
	}
	testAPI.PostJSON("/ingredients", invalidBody).
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(td.JSON(`
		{
			"error": {
				"creator_id": "creator_id id must be a positive integer",
				"name": "must not be empty"
			}
		}
		`))
}

func TestListRecipesHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	testAPI.Get("/recipes").
		Name("List all recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{
			"recipes": models.MockedRecipes,
		})

	testAPI.Get("/recipes?include=1").
		Name("List recipes including ingredient 1").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{"recipes": models.MockedRecipes[:1]})

	testAPI.Get("/recipes?exclude=1").
		Name("List recipes excluding ingredient 1").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{"recipes": models.MockedRecipes[1:]})

	testAPI.Get("/recipes?include=1&exclude=2").
		Name("List recipes including ingredient 1 and excluding 2").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{"recipes": {}})

	testAPI.Get("/recipes?include=a").
		Name("List recipes including invalid ingredient id").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`
			{
				"error": "Invalid Ingredient ID (include)"
			}
		`))

	testAPI.Get("/recipes?include=2&exclude=2").
		Name("List recipes including and excluding ingredient 2").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`
			{
				"error": "Cannot include and exclude ingredient with ID 2"
			}
		`))
}

func TestCreateRecipe(t *testing.T) {
	testAPI := newTestAPI(t)

	body := map[string]interface{}{
		"name":        "Test",
		"creator_id":  1,
		"description": "Some description",
		"ingredients": map[string]interface{}{
			"1": map[string]interface{}{
				"amount": 100,
				"unit":   "g",
			},
		},
	}

	expectedId := len(models.MockedRecipes) + 1

	testAPI.PostJSON("/recipes", body).
		Name("Create a recipe").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(td.JSON(`
		{
			"newRecipe": {
				"id": $id,
				"creator_id": 1,
				"name": "Test",
				"string": "Some description",
				"created": "$created",
				"ingredients": {
					"1": {
						"name": "Farine",
						"amount": 100,
						"unit": "g"
					},
				}
			}
		}
		`, td.Tag("id", expectedId),
			td.Tag("created", td.Between(testAPI.SentAt(), time.Now())),
		))
}

func TestListFavoritesHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	testAPI.Get("/users/1/favorites").
		Name("List user 1 favorite recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`
			{
				"favorites": [
					{
						"id": 1,
						"recipe_id": 1,
						"user_id": 1
					}
				]
			}
		`))

	testAPI.Get("/users/5/favorites").
		Name("List user 5 favorite recipes - should return empty array").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`
			{
				"favorites": []
			}
		`))

	testAPI.Get("/users/abc/favorites").
		Name("Invalid user id").
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(td.JSON(`
			{
				"error": "The requested resource could not be found"
			}
		`))
}

func TestFlagFavoriteRecipeHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	body := map[string]int{"recipe_id": 1}
	expectedNewID := len(models.MockedFavorites) + 1

	testAPI.PostJSON("/users/2/favorites", body).
		Name("Flag recipe 1 as user 2 favorite one").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(td.JSON(`
			{
				"newFavorite": {
					"id": $id,
					"recipe_id": 1,
					"user_id": 2
				}
			}
		`, td.Tag("id", expectedNewID),
		))

	invalidBody := map[string]int{"recipe_id": -1}
	testAPI.PostJSON("/users/2/favorites", invalidBody).
		Name("Flag invalid recipe -1 as user 2 favorite one").
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(td.JSON(`
			{
				"error": {
					"recipe_id": "recipe_id must be a positive integer"
				}
			}
		`))
}

func TestUnFlagFavoriteHandler(t *testing.T) {
	testAPI := newTestAPI(t)

	testAPI.Delete("/users/2/favorites/2", nil).
		Name("Delete favorite item 2 from user 2 list of favorite recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`
			{
				"message": "Deleted"
			}
		`))

	testAPI.Delete("/users/2/favorites/9", nil).
		Name("Delete out of range/unknown id").
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(td.JSON(`
			{
				"error": "The requested resource could not be found"
			}
		`))
}
