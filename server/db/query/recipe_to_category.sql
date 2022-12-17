-- name: AssignRecipeToCategory :one
INSERT INTO "RecipeCategory_Recipe" (
	id, recipe_id, recipe_category_id
) VALUES (
	$1, $2, $3
) RETURNING *;
