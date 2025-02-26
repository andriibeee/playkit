// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package pgsql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addVideoToPlaylist = `-- name: AddVideoToPlaylist :exec
INSERT INTO video(id, playlist, videoID, title, author, thumbnail, duration)
VALUES($1, $2, $3, $4, $5, $6, $7)
`

type AddVideoToPlaylistParams struct {
	ID        pgtype.UUID
	Playlist  string
	Videoid   string
	Title     string
	Author    string
	Thumbnail string
	Duration  string
}

func (q *Queries) AddVideoToPlaylist(ctx context.Context, arg AddVideoToPlaylistParams) error {
	_, err := q.db.Exec(ctx, addVideoToPlaylist,
		arg.ID,
		arg.Playlist,
		arg.Videoid,
		arg.Title,
		arg.Author,
		arg.Thumbnail,
		arg.Duration,
	)
	return err
}

const deleteVideoFromPlaylist = `-- name: DeleteVideoFromPlaylist :exec
DELETE FROM video WHERE playlist = $1 AND id = $2
`

type DeleteVideoFromPlaylistParams struct {
	Playlist string
	ID       pgtype.UUID
}

func (q *Queries) DeleteVideoFromPlaylist(ctx context.Context, arg DeleteVideoFromPlaylistParams) error {
	_, err := q.db.Exec(ctx, deleteVideoFromPlaylist, arg.Playlist, arg.ID)
	return err
}

const getPlaylist = `-- name: GetPlaylist :many
SELECT id, playlist, videoid, title, author, thumbnail, duration FROM video
WHERE playlist = $1
`

func (q *Queries) GetPlaylist(ctx context.Context, playlist string) ([]Video, error) {
	rows, err := q.db.Query(ctx, getPlaylist, playlist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Playlist,
			&i.Videoid,
			&i.Title,
			&i.Author,
			&i.Thumbnail,
			&i.Duration,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
