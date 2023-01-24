package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type RecipeIngredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"` // g | ml | ...
}

type Recipe struct {
	ID          int                      `json:"id"`
	CreatorId   int                      `json:"creator_id"`
	Name        string                   `json:"name"`
	Description string                   `json:"string"`
	Created     time.Time                `json:"created"`
	Ingredients map[int]RecipeIngredient `json:"ingredients"`
}

type RecipeModel struct {
	DB *sql.DB
}

// GetAll() returns all existing recipes
func (m *RecipeModel) GetAll(include, exclude map[int]struct{}) ([]*Recipe, error) {
	query := `
		SELECT id, name, description, creator_id, ingredients, created
		FROM recipes
	`

	if len(include) > 0 {
		query += "WHERE "
		for key := range include {
			query += fmt.Sprintf("ingredients -> '%d' IS NOT NULL AND ", key)
		}
		if len(exclude) == 0 {
			query = query[0 : len(query)-len("AND ")]
		}
	}
	if len(exclude) > 0 {
		if len(include) == 0 {
			query += "WHERE "
		}
		for key := range exclude {
			query += fmt.Sprintf("ingredients -> '%d' IS NULL AND ", key)
		}
		query = query[0 : len(query)-len("AND ")]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := []*Recipe{}
	for rows.Next() {
		var rep Recipe
		var ingredientsAsString string
		err = rows.Scan(&rep.ID, &rep.Name, &rep.Description, &rep.CreatorId, &ingredientsAsString, &rep.Created)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(ingredientsAsString), &rep.Ingredients)

		if err != nil {
			return nil, err
		}

		recipes = append(recipes, &rep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

// Insert() inserts new recipes
// and return this newly created recipe
func (m *RecipeModel) Insert(recipe *Recipe) error {
	str, err := json.Marshal(recipe.Ingredients)

	if err != nil {
		return err
	}

	query := `
		INSERT INTO recipes (name, creator_id, description, ingredients)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created
	`

	return m.DB.QueryRow(query, recipe.Name, recipe.CreatorId, recipe.Description, str).Scan(&recipe.ID, &recipe.Created)
}
