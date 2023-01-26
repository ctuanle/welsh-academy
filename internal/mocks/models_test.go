package mocks

import (
	"database/sql"
	"testing"

	"ctuanle.ovh/welsh-academy/internal/models"
	"github.com/maxatome/go-testdeep/td"
)

func TestIngredientGetAll(t *testing.T) {
	model := NewMockModels().Ingredients
	ingredients, _ := model.GetAll()

	td.Cmp(t, ingredients, MockedIngredients)
}

func TestIngredientGetById(t *testing.T) {
	model := NewMockModels().Ingredients

	// get with negative id
	ing, err := model.GetById(-1)
	td.Cmp(t, ing, td.Nil())
	td.Cmp(t, err, sql.ErrNoRows)

	// get with zero id
	ing, err = model.GetById(0)
	td.Cmp(t, ing, td.Nil())
	td.Cmp(t, err, sql.ErrNoRows)

	// get with id out of range
	ing, err = model.GetById(len(MockedIngredients) + 1)
	td.Cmp(t, ing, td.Nil())
	td.Cmp(t, err, sql.ErrNoRows)

	// get with valid id
	ing, err = model.GetById(1)
	td.Cmp(t, err, td.Nil())
	td.Cmp(t, ing, MockedIngredients[0])
}

func TestIngredientInsert(t *testing.T) {
	model := NewMockModels().Ingredients

	expectedNewId := len(MockedIngredients) + 1
	newIng := models.Ingredient{
		Name:      "Am New",
		CreatorId: 1,
	}

	_ = model.Insert(&newIng)

	td.Cmp(t, newIng.ID, expectedNewId)
}

func TestRecipeGetAll(t *testing.T) {
	model := NewMockModels().Recipes

	// without ingredient constraint
	recipes, _ := model.GetAll(nil, nil)
	td.Cmp(t, recipes, MockedRecipes)

	// with ingredient constraint
	include := map[int]struct{}{1: {}}
	exclude := map[int]struct{}{2: {}}

	recipes, _ = model.GetAll(include, nil)
	td.Cmp(t, len(recipes), 1)
	td.Cmp(t, recipes[0], MockedRecipes[0])

	recipes, _ = model.GetAll(nil, exclude)
	td.Cmp(t, len(recipes), 1)
	td.Cmp(t, recipes[0], MockedRecipes[1])

	recipes, _ = model.GetAll(include, exclude)
	td.Cmp(t, len(recipes), 0)
}

func TestRecipeInsert(t *testing.T) {
	model := NewMockModels().Recipes

	newRec := models.Recipe{
		Name:        "Test",
		CreatorId:   1,
		Description: "Test",
		Ingredients: map[int]models.RecipeIngredient{
			1: {
				Amount: 1,
				Unit:   "kg",
			},
		},
	}

	expectedNewId := len(MockedRecipes) + 1

	_ = model.Insert(&newRec)

	td.Cmp(t, newRec.ID, expectedNewId)
}

func TestFavoriteGetAll(t *testing.T) {
	model := NewMockModels().Favorites

	favorites, _ := model.GetAll(1)
	td.Cmp(t, favorites, []*models.Favorite{
		{
			ID:       1,
			RecipeId: 1,
			UserId:   1,
		},
	})

	favorites, _ = model.GetAll(4)
	td.Cmp(t, len(favorites), 0)
}

func TestFavoriteInsert(t *testing.T) {
	model := NewMockModels().Favorites

	newFav := models.Favorite{
		UserId:   4,
		RecipeId: 1,
	}

	expectedNewId := len(MockedFavorites) + 1

	model.Insert(&newFav)

	td.Cmp(t, newFav.ID, expectedNewId)
}

func TestFavoriteRemove(t *testing.T) {
	model := NewMockModels().Favorites

	err := model.Remove(-1)
	td.Cmp(t, err, sql.ErrNoRows)

	err = model.Remove(4)
	td.Cmp(t, err, sql.ErrNoRows)

	err = model.Remove(2)
	td.Cmp(t, err, td.Nil())
}
