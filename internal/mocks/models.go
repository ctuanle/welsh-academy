// This package purpose is to mimic all things models can do with mock data
// Models cannot import anything from this package but mocks can import from models
package mocks

import "ctuanle.ovh/welsh-academy/internal/models"

// NewMockModels() return mock models for testing
func NewMockModels() models.Models {
	return models.Models{
		Ingredients: MockIngredientModel{},
		Recipes:     MockRecipeModel{},
		Favorites:   MockFavoriteModel{},
	}
}
