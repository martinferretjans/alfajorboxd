-- name: CreateAlfajor :one
INSERT INTO alfajor (id, name, relleno, cobertura, precio, tapas, marca, descripcion)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAlfajor :one
SELECT *
FROM alfajor
WHERE id = $1;

-- name: ListAlfajores :many
SELECT *
FROM alfajor
ORDER BY name;

-- name: UpdateAlfajor :one
UPDATE alfajor
SET	name = $2,
	relleno = $3,
	cobertura = $4,
	precio = $5,
	tapas = $6,
	marca = $7,
	descripcion = $8
WHERE id = $1
RETURNING *;

-- name: DeleteAlfajor :exec
DELETE FROM alfajor
WHERE id = $1;
