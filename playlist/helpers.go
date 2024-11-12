package playlist

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL           = errors.New("invalid URL")
	ErrNotAnYoutubeVideoURL = errors.New("not an youtube video url")
)

func ExtractVideoID(ur string) (string, error) {
	ur = strings.TrimSpace(ur)
	if !strings.HasPrefix(ur, "http://") && !strings.HasPrefix(ur, "https://") {
		ur = "https://" + ur
	}

	u, err := url.Parse(ur)
	if err != nil {
		return "", ErrInvalidURL
	}

	if u.Host == "youtu.be" {
		return u.Path[1:], nil
	}

	if u.Host != "youtube.com" && u.Host != "www.youtube.com" {
		return "", ErrNotAnYoutubeVideoURL
	}

	q := u.Query()
	if q.Has("v") {
		return q.Get("v"), nil
	}

	return "", ErrNotAnYoutubeVideoURL
}
