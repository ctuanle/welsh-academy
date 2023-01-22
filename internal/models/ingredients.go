package models

import (
	"time"
)

type Ingredient struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Creator int       `json:"creator"`
	Created time.Time `json:"created"`
}

type IngredientModel struct {
	Ingredients []Ingredient
}

var Ingredients = []Ingredient{
	{
		ID:      1,
		Name:    "Farine",
		Creator: 1,
		Created: time.Now(),
	},
	{
		ID:      2,
		Name:    "Fromage",
		Creator: 2,
		Created: time.Now(),
	},
	{
		ID:      3,
		Name:    "Piment",
		Creator: 3,
		Created: time.Now(),
	},
}

// GetAll() returns all existing ingredients
func (m *IngredientModel) GetAll() ([]Ingredient, error) {
	return m.Ingredients, nil
}

// Insert() insert a new ingredient
// return newly created ingredient ID and an possible error
func (m *IngredientModel) Insert(name string, creator int) (*Ingredient, error) {
	newIngredient := Ingredient{
		ID:      len(m.Ingredients) + 1,
		Name:    name,
		Creator: creator,
		Created: time.Now(),
	}
	m.Ingredients = append(m.Ingredients, newIngredient)
	return &newIngredient, nil
}
