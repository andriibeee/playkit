-- name: GetPlaylist :many
SELECT * FROM video
WHERE playlist = $1;

-- name: AddVideoToPlaylist :exec
INSERT INTO video(id, playlist, videoID, title, author, thumbnail, duration)
VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteVideoFromPlaylist :exec
DELETE FROM video WHERE playlist = $1 AND id = $2;
