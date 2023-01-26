package main

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"ctuanle.ovh/welsh-academy/internal/mocks"
	"ctuanle.ovh/welsh-academy/internal/models"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

// newTestApplication() returns a new application
// for testing purpose, with mock models
func newTestApplication() *application {
	return &application{
		logger: log.New(io.Discard, "", 0),
		models: mocks.NewMockModels(),
	}
}

// newTestServer() returns a server for testing purpose
func newTestServer(app *application) *http.Server {
	return &http.Server{
		Handler: app.routes(),
	}
}

// newTestAPI() returns a new testdeep test api
func newTestAPI(t *testing.T, app *application) *tdhttp.TestAPI {
	return tdhttp.NewTestAPI(t, newTestServer(app).Handler)
}

func TestListIngredientsHandler(t *testing.T) {
	app := newTestApplication()
	testAPI := newTestAPI(t, app)
	expectedIngredients, _ := app.models.Ingredients.GetAll()

	testAPI.Get("/ingredients").
		Name("List all ingredients").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Ingredient{
			"ingredients": expectedIngredients,
		})
}

func TestCreateIngredientHandler(t *testing.T) {
	testAPI := newTestAPI(t, newTestApplication())

	var body struct {
		Name      string `json:"name"`
		CreatorId int    `json:"creator_id"`
	}

	// Test with valid input request body
	body.Name = "Test Name"
	body.CreatorId = 1
	expectedNewID := len(mocks.MockedIngredients) + 1
	testAPI.PostJSON("/ingredients", body).
		Name("Create Ingredient: (valid input)").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(map[string]*models.Ingredient{
			"ingredient": {
				ID:        expectedNewID,
				Name:      "Test Name",
				CreatorId: 1,
				Created:   testAPI.Anchor(td.Between(testAPI.SentAt(), time.Now())).(time.Time),
			},
		})

	// Test with invalid input request body
	// that fails validator
	body.Name = ""
	body.CreatorId = 0
	testAPI.PostJSON("/ingredients", body).
		Name("Create Ingredient: (input fails validator)").
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(map[string]map[string]string{
			"error": {
				"creator_id": "creator_id id must be a positive integer",
				"name":       "must not be empty",
			},
		})

	// Test with unknown user id (does not exist)
	body.Name = "Test"
	body.CreatorId = len(mocks.MockedUsers) + 5
	testAPI.PostJSON("/ingredients", body).
		Name("Create Ingredient: (creator does not exist)").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(map[string]string{
			"error": "creator id does not exist",
		})
}

func TestListRecipesHandler(t *testing.T) {
	app := newTestApplication()
	testAPI := newTestAPI(t, app)
	model := app.models.Recipes

	expectedRecipes, _ := model.GetAll(nil, nil)
	testAPI.Get("/recipes").
		Name("List all recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{
			"recipes": expectedRecipes,
		})

	expectedRecipes, _ = model.GetAll(map[int]struct{}{1: {}}, nil)
	testAPI.Get("/recipes?include=1").
		Name("List recipes including ingredient 1").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{
			"recipes": expectedRecipes,
		})

	expectedRecipes, _ = model.GetAll(nil, map[int]struct{}{1: {}})
	testAPI.Get("/recipes?exclude=1").
		Name("List recipes excluding ingredient 1").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{
			"recipes": expectedRecipes,
		})

	expectedRecipes, _ = model.GetAll(map[int]struct{}{1: {}}, map[int]struct{}{2: {}})
	testAPI.Get("/recipes?include=1&exclude=2").
		Name("List recipes including ingredient 1 and excluding 2").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Recipe{
			"recipes": expectedRecipes,
		})

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
	app := newTestApplication()
	testAPI := newTestAPI(t, app)

	type ingredientInput struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	}

	var body struct {
		Name        string                  `json:"name"`
		CreatorId   int                     `json:"creator_id"`
		Description string                  `json:"description"`
		Ingredients map[int]ingredientInput `json:"ingredients"`
	}

	// test with valid input
	body.Name = "Test"
	body.CreatorId = 1
	body.Description = "Some description"
	body.Ingredients = map[int]ingredientInput{
		1: {100, "g"},
	}

	testAPI.PostJSON("/recipes", body).
		Name("Create a recipe with valid input").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(map[string]*models.Recipe{
			"newRecipe": {
				ID:          testAPI.Anchor(td.Gt(0)).(int),
				Name:        body.Name,
				CreatorId:   body.CreatorId,
				Description: body.Description,
				Ingredients: map[int]models.RecipeIngredient{
					1: {
						Name:   "Flour",
						Amount: body.Ingredients[1].Amount,
						Unit:   body.Ingredients[1].Unit,
					},
				},
				Created: testAPI.Anchor(td.Between(testAPI.SentAt(), time.Now())).(time.Time),
			},
		})

	// test with input that fails validator
	body.Name = ""
	body.CreatorId = 0
	body.Description = ""
	body.Ingredients = map[int]ingredientInput{}

	testAPI.PostJSON("/recipes", body).
		Name("Create a recipe with valid input").
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(map[string]map[string]string{
			"error": {
				"creator_id":  "creator_id id must be a positive integer",
				"description": "recipe description can not be empty",
				"ingredients": "there must be at least one ingredient",
				"name":        "recipe name can not be empty",
			},
		})

	// test with valid input but ingredient does not exist
	body.Name = "Test"
	body.CreatorId = 1
	body.Description = "Some description"
	body.Ingredients = map[int]ingredientInput{
		50: {100, "g"},
	}

	testAPI.PostJSON("/recipes", body).
		Name("Create a recipe with valid input").
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(map[string]map[string]string{
			"error": {
				"ingredients -> id": "ingredient id 50 does not exist",
			},
		})
}

func TestListFavoritesHandler(t *testing.T) {
	app := newTestApplication()
	testAPI := newTestAPI(t, app)
	model := app.models.Favorites

	expectedFavorites, err := model.GetAll(1)
	if err != nil {
		t.Fatal(err)
	}
	testAPI.Get("/users/1/favorites").
		Name("List user 1 favorite recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Favorite{
			"favorites": expectedFavorites,
		})

	expectedFavorites, err = model.GetAll(50)
	if err != nil {
		t.Fatal(err)
	}
	testAPI.Get("/users/50/favorites").
		Name("List user 50 favorite recipes - should return empty array").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string][]*models.Favorite{
			"favorites": expectedFavorites,
		})

	testAPI.Get("/users/abc/favorites").
		Name("Invalid user id").
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(map[string]string{
			"error": "The requested resource could not be found",
		})
}

func TestFlagFavoriteRecipeHandler(t *testing.T) {
	testAPI := newTestAPI(t, newTestApplication())

	body := map[string]int{"recipe_id": 1}

	testAPI.PostJSON("/users/2/favorites", body).
		Name("Flag recipe 1 as user 2 favorite one").
		CmpStatus(http.StatusCreated).
		CmpJSONBody(map[string]*models.Favorite{
			"newFavorite": {
				ID:       testAPI.Anchor(td.NotZero(), int(0)).(int),
				RecipeId: 1,
				UserId:   2,
			},
		})

	invalidBody := map[string]int{"recipe_id": -1}
	testAPI.PostJSON("/users/2/favorites", invalidBody).
		Name("Flag invalid recipe -1 as user 2 favorite one").
		CmpStatus(http.StatusUnprocessableEntity).
		CmpJSONBody(map[string]map[string]string{
			"error": {
				"recipe_id": "recipe_id must be a positive integer",
			},
		})
}

func TestUnFlagFavoriteHandler(t *testing.T) {
	testAPI := newTestAPI(t, newTestApplication())

	testAPI.Delete("/users/2/favorites/2", nil).
		Name("Delete favorite item 2 from user 2 list of favorite recipes").
		CmpStatus(http.StatusOK).
		CmpJSONBody(map[string]string{"message": "Deleted"})

	testAPI.Delete("/users/2/favorites/9", nil).
		Name("Delete out of range/unknown id").
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(map[string]string{
			"error": "The requested resource could not be found",
		})
}
