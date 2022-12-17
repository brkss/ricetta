-- name: CreateRecipe :one

INSERT INTO "Recipe" (
	id, name , description, image, url, time, servings 
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING *;
