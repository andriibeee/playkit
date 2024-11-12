package playlist

import (
	"playlist/pgsql"
	"playlist/shared"
)

type Video struct {
	UUID      string `json:"uuid"`
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Thumbnail string `json:"thumbnail"`
	Duration  string `json:"duration"`
}

func VideoFromModel(m pgsql.Video) Video {
	return Video{
		UUID:      shared.EncodeUUID(m.ID),
		ID:        m.Videoid,
		Title:     m.Title,
		Author:    m.Author,
		Thumbnail: m.Thumbnail,
		Duration:  m.Duration,
	}
}

type TwitchEnvelope struct {
	Sender  string
	Message string
}
