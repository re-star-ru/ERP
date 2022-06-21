-- name: GetRestaritem :one
SELECT * FROM restaritems
WHERE id = $1 LIMIT 1;

-- name: ListRestaritems :many
SELECT * FROM restaritems
ORDER BY id;

-- name: CreateRestaritem :one
INSERT INTO restaritems (
	name, onceGUID
) VALUES (
			 $1, $2
		 )
RETURNING *;

-- name: DeleteRestaritem :exec
DELETE FROM restaritems
WHERE id = $1;

-- name: UpdateRestaritem :one
UPDATE restaritems
set name = $2,
	onceGUID = $3
WHERE id = $1
RETURNING *;