package shared

import "github.com/Netflix/go-env"

type Configuration struct {
	DB string `env:"DB_CONNECTION"`

	OAuth2Secret      string `env:"OAUTH2_SECRET"`
	OAuth2ID          string `env:"OAUTH2_ID"`
	OAuth2RedirectURL string `env:"OAUTH2_REDIR"`

	AuthSigningKey string `env:"AUTH_SIGNING_KEY"`

	Port string `env:"PORT"`

	YoutubeAPIKey string `env:"YOUTUBE_API_KEY"`
}

func NewConfiguration() *Configuration {
	var conf Configuration

	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
