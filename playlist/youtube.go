package playlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	ErrStatusNotOK       = errors.New("status not OK")
	ErrFailedToFindVideo = errors.New("failed to find video")
)

type YoutubeAPI struct {
	apiKey string
}

func NewYoutubeAPI(apiKey string) *YoutubeAPI {
	return &YoutubeAPI{apiKey: apiKey}
}

type VideoResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			PublishedAt  string `json:"publishedAt"`
			Thumbnails   struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

func (api *YoutubeAPI) GetVideoInfo(ctx context.Context, videoID string) (*Video, error) {
	url := fmt.Sprintf("?part=snippet,contentDetails&id=%s&key=%s", videoID, api.apiKey)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/youtube/v3/videos"+url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusNotOK
	}

	var video VideoResponse

	err = json.NewDecoder(resp.Body).Decode(&video)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %v", err) //nolint: err113
	}

	if len(video.Items) == 0 {
		return nil, ErrFailedToFindVideo
	}

	return &Video{
		ID:        videoID,
		Title:     video.Items[0].Snippet.Title,
		Author:    video.Items[0].Snippet.ChannelTitle,
		Thumbnail: video.Items[0].Snippet.Thumbnails.High.URL,
		Duration:  video.Items[0].ContentDetails.Duration,
	}, nil
}
