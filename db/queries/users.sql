-- name: CreateUser :one
INSERT INTO userdata (
  username, 
  hashed_password,
  fullname,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM userdata
WHERE username = $1 LIMIT 1;
