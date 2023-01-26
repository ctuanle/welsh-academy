package mocks

import (
	"database/sql"
	"time"

	"ctuanle.ovh/welsh-academy/internal/models"
	"github.com/lib/pq"
)

// for testing purpose
type MockIngredientModel struct{}

// GetAll() returns an array of mocked ingredients
func (m MockIngredientModel) GetAll() ([]*models.Ingredient, error) {
	return MockedIngredients, nil
}

// GetById() return a mocked ingredient
func (m MockIngredientModel) GetById(id int) (*models.Ingredient, error) {
	if id < 1 || id > len(MockedIngredients) || MockedIngredients[id-1] == nil {
		return nil, sql.ErrNoRows
	}

	return MockedIngredients[id-1], nil
}

// Insert() mocking an action of inserting an ingredient to db
func (m MockIngredientModel) Insert(ingredient *models.Ingredient) error {
	var err = pq.Error{}
	if ingredient.CreatorId > len(MockedIngredients) {
		err.Code = pq.ErrorCode("23503") // foreign_key_violation
		return &err
	}

	ingredient.ID = len(MockedIngredients) + 1
	ingredient.Created = time.Now()

	MockedIngredients = append(MockedIngredients, ingredient)

	return nil
}
