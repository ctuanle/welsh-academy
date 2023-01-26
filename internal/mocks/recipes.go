package mocks

import (
	"time"

	"ctuanle.ovh/welsh-academy/internal/models"
)

type MockRecipeModel struct{}

func (m MockRecipeModel) GetAll(include, exclude map[int]struct{}) ([]*models.Recipe, error) {
	ans := []*models.Recipe{}

	for _, rec := range MockedRecipes {
		good := true

		for in := range include {
			if _, ok := rec.Ingredients[in]; !ok {
				good = false
				break
			}
		}

		if !good {
			continue
		}

		for ex := range exclude {
			if _, ok := rec.Ingredients[ex]; ok {
				good = false
				break
			}
		}

		if good {
			ans = append(ans, rec)
		}
	}

	return ans, nil
}

func (m MockRecipeModel) Insert(recipe *models.Recipe) error {
	recipe.ID = len(MockedRecipes) + 1
	recipe.Created = time.Now()
	MockedRecipes = append(MockedRecipes, recipe)

	return nil
}
