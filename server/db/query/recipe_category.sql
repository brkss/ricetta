-- name: CreateRecipeCategory :one 
INSERT INTO "RecipeCategory" (
	id, title, image
) VALUES (
	$1, $2, $3
) RETURNING *;
