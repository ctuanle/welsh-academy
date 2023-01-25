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

type MockRecipeModel struct{}

// GetAll() returns all existing recipes
func (m RecipeModel) GetAll(include, exclude map[int]struct{}) ([]*Recipe, error) {
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
func (m RecipeModel) Insert(recipe *Recipe) error {
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

var mockedRecipes = []*Recipe{
	{
		ID:        1,
		CreatorId: 1,
		Name:      "Petits sablés",
		Ingredients: map[int]RecipeIngredient{
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
		Created:     time.Now(),
	},
	{
		ID:        2,
		CreatorId: 2,
		Name:      "Name 2",
		Ingredients: map[int]RecipeIngredient{
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
		Created:     time.Now(),
	},
}

func (m MockRecipeModel) GetAll(include, exclude map[int]struct{}) ([]*Recipe, error) {
	ans := []*Recipe{}

	for _, rec := range mockedRecipes {
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

func (m MockRecipeModel) Insert(recipe *Recipe) error {
	recipe.ID = len(mockedRecipes) + 1
	recipe.Created = time.Now()
	mockedRecipes = append(mockedRecipes, recipe)

	return nil
}
