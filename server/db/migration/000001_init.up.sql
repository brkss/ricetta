CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "userInfo" (
  "id" varchar PRIMARY KEY,
  "weight" int NOT NULL,
  "height" int NOT NULL,
  "birth" timestamp NOT NULL,
  "user_id" varchar
);

CREATE TABLE "RecipeCategory" (
  "id" varchar PRIMARY KEY,
  "title" varchar NOT NULL,
  "image" varchar,
  "active" bool DEFAULT true
);

CREATE TABLE "Recipe" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar,
  "image" varchar DEFAULT 'default.png',
  "active" bool DEFAULT true,
  "time" varchar NOT NULL,
  "url" varchar,
  "servings" int,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "RecipeCategory_Recipe" (
  "id" varchar PRIMARY KEY,
  "recipe_id" varchar,
  "recipe_category_id" varchar
);

ALTER TABLE "userInfo" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "RecipeCategory_Recipe" ADD FOREIGN KEY ("recipe_id") REFERENCES "Recipe" ("id");

ALTER TABLE "RecipeCategory_Recipe" ADD FOREIGN KEY ("recipe_category_id") REFERENCES "RecipeCategory" ("id");
