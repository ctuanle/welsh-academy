package models

import (
	"context"
	"database/sql"
	"time"
)

type Ingredient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatorId int       `json:"creator_id"`
	Created   time.Time `json:"created"`
}

type IngredientModel struct {
	DB *sql.DB
}

// for testing purpose
type MockIngredientModel struct{}

// GetAll() returns all existing ingredients
func (m IngredientModel) GetAll() ([]*Ingredient, error) {
	query := "SELECT id, name, creator_id, created FROM ingredients"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// should not last longer than 3 seconds
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := []*Ingredient{}

	for rows.Next() {
		var ing Ingredient
		err := rows.Scan(&ing.ID, &ing.Name, &ing.CreatorId, &ing.Created)

		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, &ing)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (m IngredientModel) GetById(id int) (*Ingredient, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, name, creator_id, created
		FROM ingredients
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ing Ingredient

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&ing.ID, &ing.Name, &ing.CreatorId, &ing.Created)

	if err != nil {
		return nil, err
	}

	return &ing, nil
}

// Insert() insert a new ingredient
// return newly created ingredient ID and an possible error
func (m IngredientModel) Insert(ingredient *Ingredient) error {
	query := `
		INSERT INTO ingredients (name, creator_id)
		VALUES ($1, $2)
		RETURNING id, created
	`

	return m.DB.QueryRow(query, ingredient.Name, ingredient.CreatorId).Scan(&ingredient.ID, &ingredient.Created)
}

var MockedIngredients = []*Ingredient{
	{
		ID:        1,
		Name:      "Farine",
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
		CreatorId: 3,
		Created:   time.Date(2023, 3, 24, 0, 0, 0, 0, time.UTC),
	},
}

// GetAll() returns an array of mocked ingredients
func (m MockIngredientModel) GetAll() ([]*Ingredient, error) {
	return MockedIngredients, nil
}

// GetById() return a mocked ingredient
func (m MockIngredientModel) GetById(id int) (*Ingredient, error) {
	if id < 1 || id > len(MockedIngredients) || MockedIngredients[id-1] == nil {
		return nil, sql.ErrNoRows
	}

	return MockedIngredients[id-1], nil
}

// Insert() mocking an action of inserting an ingredient to db
func (m MockIngredientModel) Insert(ingredient *Ingredient) error {
	ingredient.ID = len(MockedIngredients) + 1
	ingredient.Created = time.Now()

	MockedIngredients = append(MockedIngredients, ingredient)

	return nil
}
