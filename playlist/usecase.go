package playlist

import (
	"context"
	"encoding/json"
	"strings"
)

type PlaylistUseCase struct {
	s   *PlaylistService
	hub *Hub

	api *YoutubeAPI
}

func NewPlaylistUseCase(hub *Hub, ps *PlaylistService, api *YoutubeAPI) *PlaylistUseCase {
	return &PlaylistUseCase{
		s:   ps,
		hub: hub,
		api: api,
	}
}

func (uc *PlaylistUseCase) AddVideo(ctx context.Context, handle string, url string) error {
	url = strings.TrimSpace(url)

	u, err := ExtractVideoID(url)
	if err != nil {
		u = url
	}

	v, err := uc.api.GetVideoInfo(ctx, u)
	if err != nil {
		return err
	}

	playlist, err := uc.s.AddVideo(ctx, handle, v)
	if err != nil {
		return err
	}

	msg, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	return uc.hub.Emit(ctx, handle, msg)
}

func (uc *PlaylistUseCase) GetPlaylist(ctx context.Context, handle string) ([]Video, error) {
	return uc.s.GetVideos(ctx, handle)
}

func (uc *PlaylistUseCase) DeleteVideo(ctx context.Context, handle string, id string) error {
	err := uc.s.DeleteVideo(ctx, handle, id)
	if err != nil {
		return err
	}

	playlist, err := uc.s.GetVideos(ctx, handle)
	if err != nil {
		return err
	}

	msg, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	return uc.hub.Emit(ctx, handle, msg)
}

func (uc *PlaylistUseCase) SkipVideo(ctx context.Context, handle string) error {
	playlist, err := uc.s.PopVideo(ctx, handle)
	if err != nil {
		return err
	}

	msg, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	return uc.hub.Emit(ctx, handle, msg)
}
