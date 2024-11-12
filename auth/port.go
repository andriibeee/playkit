package auth

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vk-rv/pvx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type User struct {
	DisplayName     string `json:"display_name"`      //nolint:tagliatelle
	ProfileImageURL string `json:"profile_image_url"` //nolint:tagliatelle
}

type APIResponse struct {
	Data []User `json:"data"`
}

type AuthPort struct {
	oauth2Config *oauth2.Config
	pv4          *pvx.ProtoV4Local
	symK         *pvx.SymKey
}

func NewAuthPort(
	clientID string,
	clientSecret string,
	redirURL string,
	pv4 *pvx.ProtoV4Local,
	symK *pvx.SymKey,
) *AuthPort {
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:read:email"},
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirURL,
	}

	return &AuthPort{
		oauth2Config: oauth2Config,
		symK:         symK,
		pv4:          pv4,
	}
}

func (port *AuthPort) Login(w http.ResponseWriter, r *http.Request) {
	var tokenBytes [255]byte
	if _, err := io.ReadFull(crand.Reader, tokenBytes[:]); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	state := hex.EncodeToString(tokenBytes[:])

	http.SetCookie(w, &http.Cookie{
		Name:     "state",
		Value:    state,
		HttpOnly: true,
	})

	http.Redirect(w, r, port.oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (port *AuthPort) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	storedState, err := r.Cookie("state")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if storedState.Value != state {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	a, err := port.oauth2Config.Exchange(r.Context(), r.FormValue("code"), oauth2.SetAuthURLParam("user_info", ""))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	req.Header.Set("Client-ID", port.oauth2Config.ClientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			slog.Error("Error closing the body", slog.Any("err", cerr))
		}
	}()

	var apiResponse APIResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(apiResponse.Data) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &pvx.RegisteredClaims{
		Subject: apiResponse.Data[0].DisplayName,
	}

	token, err := port.pv4.Encrypt(port.symK, claims)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (port *AuthPort) Router(r chi.Router) {
	r.Get("/login", port.Login)
	r.Get("/callback", port.Callback)
}
