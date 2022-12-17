-- name: CreateUserInfo :one
INSERT INTO "userInfo" (
	id, weight, height, birth, user_id
) VALUES (
	$1, $2, $3, $4, $5
) RETURNING *;
