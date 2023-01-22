package models

type Favorite struct {
	ID       int `json:"id"`
	RecipeId int `json:"recipe_id"`
	UserId   int `json:"user_id"`
}

type FavoriteModel struct {
	Favorites []Favorite
}

var Favorites = []Favorite{
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

// GetAll() returns all favorite recipes of an user
func (m *FavoriteModel) GetAll(user_id int) ([]Favorite, error) {
	ans := []Favorite{}

	for _, fav := range m.Favorites {
		if fav.UserId == user_id {
			ans = append(ans, fav)
		}
	}

	return ans, nil
}

// Insert() flags recipe_id as user_id favorite
func (m *FavoriteModel) Insert(user_id, recipe_id int) (Favorite, error) {
	newFav := Favorite{
		ID:       len(m.Favorites) + 1,
		RecipeId: recipe_id,
		UserId:   user_id,
	}

	m.Favorites = append(m.Favorites, newFav)

	return newFav, nil
}

// Remove() unflags recipe_id from user_id favorite
func (m *FavoriteModel) Remove(favoriteId int) error {
	index := -1

	for i, fav := range m.Favorites {
		if fav.ID == favoriteId {
			index = i
			break
		}
	}

	if index < 0 {
		return nil
	}

	m.Favorites[index] = m.Favorites[len(m.Favorites)-1]
	m.Favorites[len(m.Favorites)-1] = Favorite{}
	m.Favorites = m.Favorites[:len(m.Favorites)-1]

	return nil
}
