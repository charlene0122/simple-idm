-- name: Create2FAInit :one
-- @two_factor_backup_codes TEXT[]
INSERT INTO login_2fa (login_id, two_factor_secret, two_factor_enabled, two_factor_type, two_factor_backup_codes)
VALUES ($1, $2, FALSE, $3, @two_factor_backup_codes::TEXT[])
RETURNING id;

-- name: Enable2FA :exec
UPDATE login_2fa
SET two_factor_enabled = TRUE,
    updated_at = now() AT TIME ZONE 'utc'
WHERE login_id = $1
AND two_factor_type = $2
AND deleted_at IS NULL;

-- name: Disable2FA :exec
UPDATE login_2fa
SET two_factor_enabled = FALSE,
    updated_at = now() AT TIME ZONE 'utc'
WHERE login_id = $1
AND two_factor_type = $2
AND deleted_at IS NULL;

-- name: Delete2FA :exec
UPDATE login_2fa
SET deleted_at = now() AT TIME ZONE 'utc'
WHERE login_id = $1
AND two_factor_type = $2
AND deleted_at IS NULL;

-- name: Get2FAByLoginId :one
SELECT id, login_id, two_factor_secret, two_factor_enabled
FROM login_2fa
WHERE login_id = $1
AND two_factor_type = $2
AND deleted_at IS NULL;

-- name: FindEnabledTwoFAs :many
SELECT id, login_id, two_factor_type, two_factor_enabled, created_at
FROM login_2fa
WHERE login_id = $1
AND two_factor_enabled = TRUE
AND deleted_at IS NULL;

