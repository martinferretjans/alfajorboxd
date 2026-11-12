-- name: CreateAlfajor :one
INSERT INTO alfajor (id, name, relleno, cobertura, precio, tapas, marca, descripcion, imagen)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetAlfajor :one
SELECT *
FROM alfajor
WHERE id = $1;

-- name: ListAlfajores :many
SELECT *
FROM alfajor
ORDER BY name;

-- name: ListAlfajoresPopulares :many
SELECT a.*, COALESCE(ROUND(AVG(r.calificacion) / 2.0, 2), 0)::float AS promedio, COUNT(r.id_user) AS cantidad
FROM alfajor a
LEFT JOIN review r ON r.id_alfajor = a.id
GROUP BY a.id
ORDER BY promedio DESC, a.name;

-- name: UpdateAlfajor :one
UPDATE alfajor
SET	name = $2,
	relleno = $3,
	cobertura = $4,
	precio = $5,
	tapas = $6,
	marca = $7,
	descripcion = $8,
	imagen = $9
WHERE id = $1
RETURNING *;

-- name: DeleteAlfajor :exec
DELETE FROM alfajor
WHERE id = $1;

-- name: CreateUsuario :one
INSERT INTO usuario (name, email)
VALUES ($1, $2)
RETURNING *;

-- name: GetUsuario :one
SELECT *
FROM usuario
WHERE id = $1;

-- name: ListUsuarios :many
SELECT *
FROM usuario
ORDER BY name;

-- name: DeleteUsuario :exec
DELETE FROM usuario
WHERE id = $1;

-- name: CreateReview :one
INSERT INTO review (id_alfajor, id_user, calificacion, comentario)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReview :one
SELECT *
FROM review
WHERE id_alfajor = $1 AND id_user = $2;

-- name: ListReviewsByAlfajor :many
SELECT r.*, u.name AS usuario
FROM review r
JOIN usuario u ON u.id = r.id_user
WHERE r.id_alfajor = $1
ORDER BY r.fecha DESC, u.name;

-- name: ListUltimasReviews :many
SELECT r.*, u.name AS usuario, a.name AS alfajor, a.imagen
FROM review r
JOIN usuario u ON u.id = r.id_user
JOIN alfajor a ON a.id = r.id_alfajor
ORDER BY r.fecha DESC, a.name
LIMIT $1;

-- name: GetPromedioAlfajor :one
SELECT COALESCE(ROUND(AVG(calificacion) / 2.0, 2), 0)::float AS promedio
FROM review
WHERE id_alfajor = $1;

-- name: DeleteReview :exec
DELETE FROM review
WHERE id_alfajor = $1 AND id_user = $2;
