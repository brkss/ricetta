-- name: CreateRecipe :one

INSERT INTO "Recipe" (
	id, name , description, image, url, time, servings 
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: GetRecipes :many
SELECT * FROM "Recipe"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetRecipeByCategory :many

SELECT * from "Recipe" as r 
LEFT JOIN "RecipeCategory_Recipe" as rcr  
ON rcr.recipe_id = r.id 
WHERE rcr.recipe_category_id = $1
AND rcr.recipe_id is NULL;