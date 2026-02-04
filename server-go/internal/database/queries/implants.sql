-- name: GetImplantByID :one
SELECT * FROM implants WHERE implant_id = $1;

-- name: GetImplantByPublicKey :one
SELECT * FROM implants WHERE implant_pk = $1;

-- name: CreateImplant :one
INSERT INTO implants (
    implant_id, hostname, ip, os, arch, username, server_sk, implant_pk
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateImplantLastSeen :exec
UPDATE implants
SET last_seen = NOW()
WHERE implant_id = $1;

-- name: DeleteImplant :exec
DELETE FROM implants WHERE implant_id = $1;

-- name: ListImplants :many
SELECT * FROM implants ORDER BY last_seen DESC;

-- name: CountImplants :one
SELECT COUNT(*) FROM implants;

-- name: CreateImplantAlias :one
INSERT INTO implant_aliases (implant_id, alias)
VALUES ($1, $2)
RETURNING *;

-- name: GetAliasByName :one
SELECT * FROM implant_aliases WHERE alias = $1;

-- name: GetAliasesByImplantID :many
SELECT * FROM implant_aliases WHERE implant_id = $1;

-- name: DeleteAlias :exec
DELETE FROM implant_aliases WHERE alias = $1;

-- name: ResolveImplantIdentifier :one
SELECT i.* FROM implants i
LEFT JOIN implant_aliases a ON i.implant_id = a.implant_id
WHERE i.implant_id = $1 OR a.alias = $1
LIMIT 1;
