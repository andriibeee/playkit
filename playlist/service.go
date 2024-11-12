package playlist

import (
	"context"
	"log/slog"
	"playlist/pgsql"
	"playlist/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistService struct {
	pool *pgxpool.Pool
}

func NewPlaylistService(pool *pgxpool.Pool) *PlaylistService {
	return &PlaylistService{
		pool: pool,
	}
}

func (ps *PlaylistService) GetVideos(ctx context.Context, handle string) ([]Video, error) {
	videos := make([]Video, 0)

	c, err := ps.pool.Acquire(ctx)
	if err != nil {
		return videos, err
	}

	defer c.Release()

	p, err := pgsql.New(c).GetPlaylist(ctx, handle)
	if err != nil {
		return videos, err
	}

	for _, v := range p {
		videos = append(videos, VideoFromModel(v))
	}

	return videos, nil
}

func (ps *PlaylistService) AddVideo(ctx context.Context, handle string, video *Video) (v []Video, err error) {
	v = make([]Video, 0)

	c, err := ps.pool.Acquire(ctx)
	if err != nil {
		return v, err
	}

	defer c.Release()

	tx, err := c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return v, err
	}

	defer func() {
		if err != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil {
				slog.Error("Rollback failed", slog.Any("error", rErr))
			}
		} else {
			cErr := tx.Commit(ctx)
			if cErr != nil {
				slog.Error("Commit failed", slog.Any("error", cErr))
			}
		}
	}()

	q := pgsql.New(tx)

	err = q.AddVideoToPlaylist(ctx, pgsql.AddVideoToPlaylistParams{
		ID:        shared.NewUUID(),
		Playlist:  handle,
		Videoid:   video.ID,
		Title:     video.Title,
		Author:    video.Author,
		Duration:  video.Duration,
		Thumbnail: video.Thumbnail,
	})
	if err != nil {
		return v, err
	}

	pl, err := q.GetPlaylist(ctx, handle)
	if err != nil {
		return v, err
	}

	for _, vd := range pl {
		v = append(v, VideoFromModel(vd))
	}

	return v, err
}

func (ps *PlaylistService) PopVideo(ctx context.Context, handle string) (v []Video, err error) {
	v = make([]Video, 0)

	c, err := ps.pool.Acquire(ctx)
	if err != nil {
		return v, err
	}

	defer c.Release()

	tx, err := c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return v, err
	}

	defer func() {
		if err != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil {
				slog.Error("Rollback failed", slog.Any("error", rErr))
			}
		} else {
			cErr := tx.Commit(ctx)
			if cErr != nil {
				slog.Error("Commit failed", slog.Any("error", cErr))
			}
		}
	}()

	q := pgsql.New(tx)

	pl, err := q.GetPlaylist(ctx, handle)
	if err != nil {
		return v, err
	}

	for _, vd := range pl {
		v = append(v, VideoFromModel(vd))
	}

	if len(pl) > 0 {
		err = q.DeleteVideoFromPlaylist(ctx, pgsql.DeleteVideoFromPlaylistParams{
			ID:       pl[0].ID,
			Playlist: handle,
		})
	}

	return v, err
}

func (ps *PlaylistService) DeleteVideo(ctx context.Context, handle string, id string) (err error) {
	c, err := ps.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer c.Release()

	tx, err := c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil {
				slog.Error("Rollback failed", slog.Any("error", rErr))
			}
		} else {
			cErr := tx.Commit(ctx)
			if cErr != nil {
				slog.Error("Commit failed:", slog.Any("error", cErr))
			}
		}
	}()

	uuid := pgtype.UUID{}

	err = uuid.Scan(id)
	if err != nil {
		return err
	}

	q := pgsql.New(tx)
	err = q.DeleteVideoFromPlaylist(ctx, pgsql.DeleteVideoFromPlaylistParams{
		ID:       uuid,
		Playlist: handle,
	})

	return err
}
