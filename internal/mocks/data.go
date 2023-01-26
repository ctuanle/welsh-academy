package mocks

import (
	"time"

	"ctuanle.ovh/welsh-academy/internal/models"
)

type User struct {
	ID       int
	UserName string
	Role     string
}

var MockedUsers = []User{
	{1, "alice", "expert"},
	{2, "bob", "expert"},
	{3, "lucy", "normal"},
}

var MockedIngredients = []*models.Ingredient{
	{
		ID:        1,
		Name:      "Flour",
		CreatorId: 1,
		Created:   time.Date(2023, 1, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:        2,
		Name:      "Fromage",
		CreatorId: 2,
		Created:   time.Date(2023, 2, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:        3,
		Name:      "Piment",
		CreatorId: 1,
		Created:   time.Date(2023, 3, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:        4,
		Name:      "Crème",
		CreatorId: 2,
		Created:   time.Date(2023, 3, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:        5,
		Name:      "Lait",
		CreatorId: 1,
		Created:   time.Date(2023, 3, 24, 0, 0, 0, 0, time.UTC),
	},
}

var MockedRecipes = []*models.Recipe{
	{
		ID:        1,
		CreatorId: 1,
		Name:      "Petits sablés",
		Ingredients: map[int]models.RecipeIngredient{
			1: {
				Name:   "Fromage",
				Amount: 100,
				Unit:   "g",
			},
			2: {
				Name:   "Piment",
				Amount: 50,
				Unit:   "g",
			},
			3: {
				Name:   "Crème",
				Amount: 100,
				Unit:   "g",
			},
		},
		Description: "This is a simple description",
		Created:     time.Date(2023, 1, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:        2,
		CreatorId: 2,
		Name:      "Name 2",
		Ingredients: map[int]models.RecipeIngredient{
			3: {
				Name:   "Fromage",
				Amount: 100,
				Unit:   "g",
			},
			4: {
				Name:   "Piment",
				Amount: 50,
				Unit:   "g",
			},
			5: {
				Name:   "Crème",
				Amount: 100,
				Unit:   "g",
			},
		},
		Description: "This is a simple description",
		Created:     time.Date(2023, 1, 24, 0, 0, 0, 0, time.UTC),
	},
}

var MockedFavorites = []*models.Favorite{
	{
		ID:       1,
		RecipeId: 1,
		UserId:   1,
	},
	{
		ID:       2,
		RecipeId: 1,
		UserId:   2,
	},
	{
		ID:       3,
		RecipeId: 2,
		UserId:   3,
	},
}
