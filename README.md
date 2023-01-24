# Welsh Academy

**Welsh Academy** is an application dedicated to provide recipes to cheddar lovers around the world. This a backend allowing cheddar **experts** to:

- Create ingredients
- Create recipes of meals using the previously created ingredients

A **user** should be able to enjoy the recipes by using the API to:

- list all existing ingredients
- list all possible recipes (with or without ingredient constraints)
- flag/unflag recipes as his favorite ones
- list his favorite recipes

## Table of contents

- [Tables of endpoints](./README.md#table-of-endpoints)
- [Entity/Model design](./README.md#entitymodel-design)
  - [User](./README.md#user)
  - [Ingredient](./README.md#ingredient)
  - [Recipe](./README.md#recipe)
  - [Favorite Recipe](./README.md#favorite-recipe)
- [Usage](./README.md#usage)
  - [/ingredients](./README.md#ingredients)
  - [/recipes](./README.md#recipes)
  - [/users/:uid/favorites](./README.md#usersuidfavorites)

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
            <td>/recipes?include=1,2&exclude=3,4</td>
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
            <td>/users/:uid/favorites/:fid</td>
            <td>unflagFavoriteRecipe</td>
            <td>Unflag/Remove favorite with id "fid"</td>
        </tr>
    </tbody>
</table>

## Entity/Model Design

### User

User is supposed to be managed by another api, but here I would like to have a simple structure of an user

<table>
    <tbody>
        <tr>
            <td><b>Field</b></td>
            <td>ID</td>
            <td>Username</td>
            <td>Role</td>
        </tr>
        <tr>
            <td><b>Type</b></td>
            <td>Integer</td>
            <td>String</td>
            <td>Expert/User</td>
        </tr>
    </tbody>
</table>

### Ingredient

<table>
    <tbody>
        <tr>
            <td><b>Field</b></td>
            <td>ID</td>
            <td>Name</td>
            <td>CreatorId</td>
            <td>Created</td>
        </tr>
        <tr>
            <td><b>Type</b></td>
            <td>Integer</td>
            <td>String</td>
            <td>Integer</td>
            <td>Time</td>
        </tr>
    </tbody>
</table>

### Recipe

**SubType(RecipeIngredient)**

<table>
    <tbody>
        <tr>
            <td><b>Field</b></td>
            <td>ID(IngredientID)</td>
            <td>Amount</td>
            <td>Unit</td>
        </tr>
        <tr>
            <td><b>Type</b></td>
            <td>Integer</td>
            <td>Float</td>
            <td>String (ml/g/l/kg/...)</td>
        </tr>
    </tbody>
</table>

**Recipe**

<table>
    <tbody>
        <tr>
            <td><b>Field</b></td>
            <td>ID</td>
            <td>CreatorId</td>
            <td>Name</td>
            <td>Ingredients</td>
            <td>Description</td>
            <td>Created</td>
        </tr>
        <tr>
            <td><b>Type</b></td>
            <td>Integer</td>
            <td>Integer</td>
            <td>String</td>
            <td>[]RecipeIngredient</td>
            <td>String</td>
            <td>Created</td>
        </tr>
    </tbody>
</table>

### Favorite Recipe

<table>
    <tbody>
        <tr>
            <td><b>Field</b></td>
            <td>ID</td>
            <td>UserId</td>
            <td>RecipeId</td>
        </tr>
        <tr>
            <td><b>Type</b></td>
            <td>Integer</td>
            <td>Integer</td>
            <td>Integer</td>
        </tr>
    </tbody>
</table>

## Usage

### **/ingredients**

**GET /ingredients**

**POST /ingredient**
with body

```json
{
  "name": "Tomato",
  "creator_id": 3
}
```

### **/recipes**

**GET /recipes** </br>
**GET /recipes?include=1** </br>
**GET /recipes?include=1,2** </br>
**GET /recipes?exclude=1,2** </br>
**GET /recipes?include=1,3&exclude=2** </br>

**POST /recipes** with body

```json
{
  "name": "Name of recipe",
  "description": "Some details, steps, bla... description",
  "creator_id": 3,
  "ingredients": {
    "1": {
      "amount": 0.5,
      "unit": "bow"
    },
    "2": {
      "amount": 1,
      "unit": "spoon"
    },
    "3": {
      "amount": 100,
      "unit": "g"
    }
  }
}
```

### **/users/:uid/favorites**

**GET /users/1/favorites** </br>
**GET /users/2/favorites** </br>

**POST /users/1/favorites** </br>
with body

```json
{
  "recipe_id": 1
}
```

**DELETE /users/1/favorites/1**

## Docker Container

We need two environnement variables : POSTGRES_PASSWORD and POSTGRES_DNS.

- POSTGRES_PASSWORD to set password for default postgres user
- POSTGRES_DNS for app to connect

For example:

- POSTGRES_PASSWORD=password
- POSTGRES_DNS=postgres://postgres:password@postgresql/welsh?sslmode=disable

```shell
cd ./path/to/root/of/project
export POSTGRES_PASSWORD=password
export POSTGRES_DNS=postgres://postgres:password@postgresql/welsh?sslmode=disable
docker compose build
docker compose up
```

For the moment, there might be an error that the app cannot connect to service postgresql, stop and run <code>docker compose up</code> will fix that.
