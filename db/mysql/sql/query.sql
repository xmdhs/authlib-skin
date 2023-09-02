-- name: GetUser :one
SELECT  *
FROM user
WHERE id = ?
LIMIT 1;
-- name: ListUser :many 
SELECT  *
FROM user
ORDER BY reg_time;
-- name: CreateUser :execresult 
 REPLACE INTO user ( id, email, password, salt, state, reg_time ) VALUES (?, ?, ?, ?, ?, ?);
-- name: DeleteUser :exec 
 DELETE
FROM user
WHERE id = ?;
-- name: CreateUserProfile :execresult 
 REPLACE INTO `user_profile` (`user_id`, `name`, `uuid`) VALUES (?, ?, ?);
-- name: GetUserByEmail :one 
SELECT  *
FROM user
WHERE email = ?
LIMIT 1 for update;
-- name: GetUserProfileByName :one 
SELECT  *
FROM `user_profile`
WHERE `name` = ?
LIMIT 1 for update;