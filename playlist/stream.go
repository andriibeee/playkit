package playlist

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/vk-rv/pvx"
)

type Streamer struct {
	h    *Hub
	uc   *PlaylistUseCase
	symK *pvx.SymKey
	pv4  *pvx.ProtoV4Local
}

func NewStreamer(
	h *Hub,
	uc *PlaylistUseCase,
	symK *pvx.SymKey,
	pv4 *pvx.ProtoV4Local,
) *Streamer {
	return &Streamer{
		h:    h,
		uc:   uc,
		symK: symK,
		pv4:  pv4,
	}
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("token") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	cc := pvx.RegisteredClaims{}

	err := s.pv4.
		Decrypt(token, s.symK).
		ScanClaims(&cc)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	handle := cc.Subject

	b := NewBot(handle)
	b.OnMessage(ctx, func(sender, msg string) {
		message := strings.TrimSpace(msg)
		if strings.HasPrefix(message, "!play ") {
			message = strings.Replace(message, "!play ", "", 1)
			message = strings.TrimSpace(message)

			if message != "" {
				err := s.uc.AddVideo(ctx, handle, message)
				if err != nil {
					slog.Error("failed to add video", slog.Any("error", err))
				}
			}
		}
	})

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("failed to hijack", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer func(c *websocket.Conn) {
		cerr := c.CloseNow()
		if cerr != nil {
			slog.Error("failed to close", slog.Any("error", err))
		}
	}(c)

	s.h.Register(handle, c)
	defer s.h.Delete(handle)

	err = b.Start(ctx)
	if err != nil {
		slog.Error("Error starting bot", slog.Any("error", err))

		return
	}

	for {
		<-ctx.Done()

		return
	}
}
