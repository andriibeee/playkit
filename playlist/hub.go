package playlist

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/coder/websocket"
)

var ErrClientNotFound = errors.New("client not found")

type Client struct {
	Handle string
	Conn   *websocket.Conn
}

type Hub struct {
	clients sync.Map
}

func NewHub() *Hub {
	return &Hub{
		clients: sync.Map{},
	}
}

func (h *Hub) Register(handle string, conn *websocket.Conn) {
	existing, ok := h.clients.Load(handle)
	if ok {
		err := existing.(*Client).Conn.CloseNow()
		if err != nil {
			slog.Error("failed to close existing connection", slog.Any(
				"handle", handle,
			))
		}

		h.clients.Delete(handle)
	}

	h.clients.Store(handle, &Client{
		Handle: handle,
		Conn:   conn,
	})
}

func (h *Hub) Delete(handle string) {
	h.clients.Delete(handle)
}

func (h *Hub) Emit(ctx context.Context, handle string, message []byte) error {
	val, ok := h.clients.Load(handle)
	if !ok {
		return ErrClientNotFound
	}

	cli, ok := val.(*Client)
	if !ok {
		return ErrClientNotFound
	}

	conn := cli.Conn

	return conn.Write(ctx, websocket.MessageText, message)
}
