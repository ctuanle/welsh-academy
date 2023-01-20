# Welsh Academy

**Welsh Academy** is an application dedicated to provide recipes to cheddar lovers around the world. This a backend allowing cheddar **experts** to:

- Create ingredients
- Create recipes of meals using the previously created ingredients

A **user** should be able to enjoy the recipes by using the API to:

- list all existing ingredients
- list all possible recipes (with or without ingredient constraints)
- flag/unflag recipes as his favorite ones
- list his favorite recipes

## Table of endpoints

<table>
    <thead>
        <tr>
            <th>Method</th>
            <th>URL Pattern</th>
            <th>Handler</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>GET</td>
            <td>/</td>
            <td></td>
            <td>Display a welcome message</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/ingredients</td>
            <td>listIngredients</td>
            <td>List all existing ingredients</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/ingredients</td>
            <td>createIngredient</td>
            <td>Create a new ingredient (expert only)</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/recipes</td>
            <td>listRecipes</td>
            <td>List all existing recipes</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/recipes</td>
            <td>createRecipe</td>
            <td>Create a new recipe (expert only)</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/users/:uid/favorites</td>
            <td>listFavorites</td>
            <td>List all favorite recipes of an user with id "uid"</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/users/:uid/favorites</td>
            <td>flagFavoriteRecipe</td>
            <td>Flag/Add a recipe as user favorite one</td>
        </tr>
        <tr>
            <td>DELETE</td>
            <td>/users/:uid/favorites/:rid</td>
            <td>unflagFavoriteRecipe</td>
            <td>Unflag/Remove recipe (with id "rid") from user (with id "uid") favorite recipes list</td>
        </tr>
    </tbody>
</table>
