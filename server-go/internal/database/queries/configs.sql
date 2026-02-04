-- name: GetImplantConfig :one
SELECT * FROM implant_configs WHERE implant_id = $1;

-- name: CreateImplantConfig :one
INSERT INTO implant_configs (
    implant_id, cookie, kill_date, user_agent, auto_self_destruct,
    sleep_time, jitter, max_retries, retry_wait, retry_jitter,
    enc_key, tailoring_hash_function, tailoring_hash_rounds, tailoring_hashes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: UpdateImplantConfig :exec
UPDATE implant_configs
SET
    cookie = COALESCE(sqlc.narg('cookie'), cookie),
    kill_date = COALESCE(sqlc.narg('kill_date'), kill_date),
    user_agent = COALESCE(sqlc.narg('user_agent'), user_agent),
    auto_self_destruct = COALESCE(sqlc.narg('auto_self_destruct'), auto_self_destruct),
    sleep_time = COALESCE(sqlc.narg('sleep_time'), sleep_time),
    jitter = COALESCE(sqlc.narg('jitter'), jitter),
    max_retries = COALESCE(sqlc.narg('max_retries'), max_retries),
    retry_wait = COALESCE(sqlc.narg('retry_wait'), retry_wait),
    retry_jitter = COALESCE(sqlc.narg('retry_jitter'), retry_jitter),
    tailoring_hash_function = COALESCE(sqlc.narg('tailoring_hash_function'), tailoring_hash_function),
    tailoring_hash_rounds = COALESCE(sqlc.narg('tailoring_hash_rounds'), tailoring_hash_rounds),
    tailoring_hashes = COALESCE(sqlc.narg('tailoring_hashes'), tailoring_hashes)
WHERE implant_id = sqlc.arg('implant_id');

-- name: DeleteImplantConfig :exec
DELETE FROM implant_configs WHERE implant_id = $1;
