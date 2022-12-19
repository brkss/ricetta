-- name: CreateRecipeCategory :one 
INSERT INTO "RecipeCategory" (
	id, title, image
) VALUES (
	$1, $2, $3
) RETURNING *;


-- name: GetCategories :many
SELECT * FROM "RecipeCategory"
ORDER BY id
LIMIT $1
OFFSET $2;