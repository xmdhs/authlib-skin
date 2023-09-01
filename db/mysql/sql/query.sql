-- name: GetUser :one
SELECT
  *
FROM
  user
WHERE
  id = ?
LIMIT
  1;

-- name: ListUser :many
SELECT
  *
FROM
  user
ORDER BY
  reg_time;

-- name: CreateUser :execresult
INSERT INTO
  user (
    id,
    email,
    password,
    salt,
    disabled,
    admin,
    reg_time
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM
  user
WHERE
  id = ?;