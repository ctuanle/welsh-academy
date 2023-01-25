# Welsh Academy

**Welsh Academy** est une application destinée à fournir des recettes aux amateurs de cheddar du monde entier. Il s'agit d'un backend permettant aux **experts** du cheddar de:

- Créer des ingrédients
- Créer des recettes de repas en utilisant les ingrédients créés précédemment

Un **utilisateur** devrait pouvoir profiter des recettes en utilisant l'API pour :

- lister tous les ingrédients existants
- lister toutes les recettes possibles (avec ou sans contrainte d'ingrédient)
- marquer/désigner ses recettes préférées
- lister ses recettes préférées

<img src="./architect.jpg" alt="Architect">

## Table des matières

- [Tables of endpoints](#table-of-endpoints)
- [Entity/Model design](#entitymodel-design)
  - [User](#user)
  - [Ingredient](#ingredient)
  - [Recipe](#recipe)
  - [Favorite Recipe](#favorite-recipe)
- [API Usage](#api-usage)
  - [/ingredients](#ingredients)
  - [/recipes](#recipes)
  - [/users/:uid/favorites](#usersuidfavorites)
- [Running](#running)
  - [With docker](#with-docker)
  - [Without docker](#without-docker)
- [Testing](#testing)

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
            <td>Afficher un message de bienvenue</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/ingredients</td>
            <td>listIngredients</td>
            <td>Liste de tous les ingrédients existants</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/ingredients</td>
            <td>createIngredient</td>
            <td>Créer un nouvel ingrédient (expert uniquement)</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/recipes?include=1,2&exclude=3,4</td>
            <td>listRecipes</td>
            <td>Lister toutes les recettes existantes</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/recipes</td>
            <td>createRecipe</td>
            <td>Créer une nouvelle recette (expert uniquement)</td>
        </tr>
        <tr>
            <td>GET</td>
            <td>/users/:uid/favorites</td>
            <td>listFavorites</td>
            <td>Liste de toutes les recettes préférées d'un utilisateur avec l'id "uid".</td>
        </tr>
        <tr>
            <td>POST</td>
            <td>/users/:uid/favorites</td>
            <td>flagFavoriteRecipe</td>
            <td>Marquer/Ajouter une recette comme favorite de l'utilisateur</td>
        </tr>
        <tr>
            <td>DELETE</td>
            <td>/users/:uid/favorites/:fid</td>
            <td>unflagFavoriteRecipe</td>
            <td>Déclarer/supprimer un favori avec l'identifiant "fid".</td>
        </tr>
    </tbody>
</table>

## Entity/Model Design

### User

L'utilisateur est censé être géré par une autre api, mais ici je voudrais avoir une structure simple d'un utilisateur

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

## API Usage

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

## Running

### **With docker**

Nous avons besoin de deux variables d'environnement : POSTGRES_PASSWORD et POSTGRES_DNS.

- POSTGRES_PASSWORD pour définir le mot de passe de l'utilisateur postgres par défaut.
- POSTGRES_DNS pour que l'application puisse se connecter.

Par exemple :

- POSTGRES_PASSWORD=password
- POSTGRES_DNS=postgres://postgres:password@postgresql/welsh?sslmode=disable

```shell
cd ./path/to/root/of/project
export POSTGRES_PASSWORD=password
export POSTGRES_DNS=postgres://postgres:password@postgresql/welsh?sslmode=disable
docker compose build
docker compose up
```

L'api sera disponible sur localhost:8000

Pour le moment, il pourrait y avoir une erreur que l'application ne peut pas se connecter au service postgresql, arrêter et exécuter <code>docker compose up</code> à nouveau corrigera cela.

### **Without docker**

Premières migrations (uniquement pour la première exécution de l'application)

```shell
go run ./cmd/cli/migrate -dns=your_postgres_dns -up
```

Lancer l'application

```shell
go run ./cmd/api -db-dns=your__postgres_dns -port=your_port
```

## Testing

Test des gestionnaires d'API

```shell
CGO_ENABLED=0 go test ./cmd/api
```

Test des modèles de BD (mock)

```shell
go test ./internal/models/
```
