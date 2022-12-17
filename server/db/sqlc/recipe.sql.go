// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: recipe.sql

package db

import (
	"context"
)

const createRecipe = `-- name: CreateRecipe :one

INSERT INTO "Recipe" (
	id, name , description, image, url, time, servings 
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING id, name, description, image, active, time, url, servings, created_at
`

type CreateRecipeParams struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Url         string `json:"url"`
	Time        string `json:"time"`
	Servings    int32  `json:"servings"`
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error) {
	row := q.db.QueryRowContext(ctx, createRecipe,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Image,
		arg.Url,
		arg.Time,
		arg.Servings,
	)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Image,
		&i.Active,
		&i.Time,
		&i.Url,
		&i.Servings,
		&i.CreatedAt,
	)
	return i, err
}
