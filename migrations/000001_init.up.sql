-- migrate -path=./migrations -database=postgres://welsh:azerty@localhost/welsh up
-- migrate -path=./migrations -database=postgres://welsh:azerty@localhost/welsh down
-- migrate -path=./migrations -database=postgres://welsh:azerty@localhost/welsh force 1

CREATE TYPE role as ENUM ('normal', 'expert');

CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  username text NOT NULL,
  role role DEFAULT 'normal'
);

INSERT INTO users (username, role) VALUES ('bob', 'expert');
INSERT INTO users (username, role) VALUES ('alice', 'expert');
INSERT INTO users (username, role) VALUES ('lucy', 'expert');
INSERT INTO users (username) VALUES ('gopher');
INSERT INTO users (username) VALUES ('mountain');
INSERT INTO users (username) VALUES ('cloud');
INSERT INTO users (username) VALUES ('ocean');


CREATE TABLE IF NOT EXISTS ingredients (
  id serial PRIMARY KEY,
  created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  creator_id int REFERENCES users(id)
);

INSERT INTO ingredients (name, creator_id) VALUES ('Farine', 1);
INSERT INTO ingredients (name, creator_id) VALUES ('Huile d Olive', 2);
INSERT INTO ingredients (name, creator_id) VALUES ('Piment', 2);
INSERT INTO ingredients (name, creator_id) VALUES ('Fromage', 3);
INSERT INTO ingredients (name, creator_id) VALUES ('Chip mexicaine', 1);
INSERT INTO ingredients (name, creator_id) VALUES ('Haricot rouge', 2);
INSERT INTO ingredients (name, creator_id) VALUES ('Mai', 3);
INSERT INTO ingredients (name, creator_id) VALUES ('Oignon rouge', 3);
INSERT INTO ingredients (name, creator_id) VALUES ('Sauce salsa', 2);
INSERT INTO ingredients (name, creator_id) VALUES ('Cheddar à fondre', 2);
INSERT INTO ingredients (name, creator_id) VALUES ('Crème fraiche', 1);
INSERT INTO ingredients (name, creator_id) VALUES ('Guacamole', 1);

CREATE TABLE IF NOT EXISTS recipes (
  id serial PRIMARY KEY,
  created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  creator_id INT REFERENCES users(id),
  description text NOT NULL
);

CREATE TABLE IF NOT EXISTS recipeingredient (
  recipe_id int REFERENCES recipes(id),
  ingredient_id int REFERENCES ingredients(id),
  amount decimal NOT NULL,
  unit text NOT NULL,
  CONSTRAINT recipe_ingredient_pkey PRIMARY KEY (recipe_id, ingredient_id)
);

INSERT INTO recipes (name, creator_id, description) VALUES ('Les nachos au four', 1, 'Ici une détailée description');
INSERT INTO recipeingredient (recipe_id, ingredient_id, amount, unit) VALUES (1, 5, 1, 'paquet');
INSERT INTO recipeingredient (recipe_id, ingredient_id, amount, unit) VALUES (1, 6, 1, 'petit boite');
INSERT INTO recipeingredient (recipe_id, ingredient_id, amount, unit) VALUES (1, 7, 1, 'petit boite');
INSERT INTO recipeingredient (recipe_id, ingredient_id, amount, unit) VALUES (1, 10, 0.5, 'tranche');

CREATE TABLE IF NOT EXISTS favorites (
  id serial PRIMARY KEY,
  recipe_id int REFERENCES recipes(id),
  user_id int REFERENCES users(id)
);

INSERT INTO favorites (recipe_id, user_id) VALUES (1, 1);
INSERT INTO favorites (recipe_id, user_id) VALUES (1, 1);