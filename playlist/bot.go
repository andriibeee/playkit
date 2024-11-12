package playlist

import (
	"context"
	"log/slog"

	"github.com/coder/websocket"
	"gopkg.in/irc.v4"
)

type Bot struct {
	handle string
	ch     chan TwitchEnvelope
}

func NewBot(handle string) *Bot {
	ch := make(chan TwitchEnvelope, 256)

	return &Bot{
		handle: handle,
		ch:     ch,
	}
}

func (bot *Bot) OnMessage(ctx context.Context, cb func(sender string, message string)) {
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case msg := <-bot.ch:
				cb(msg.Sender, msg.Message)
			}
		}
	}()
}

func (bot *Bot) Start(ctx context.Context) error {
	c, _, err := websocket.Dial(ctx, "wss://irc-ws.chat.twitch.tv/", nil) //nolint: bodyclose
	if err != nil {
		return err
	}

	initMessages := []string{
		"CAP REQ :twitch.tv/tags twitch.tv/commands",
		"PASS SCHMOOPIIE",
		"NICK justinfan67638",
		"USER justinfan67638 8 * :justinfan67638",
	}

	go func() {
		for _, msg := range initMessages {
			wErr := c.Write(ctx, websocket.MessageText, []byte(msg))
			if wErr != nil {
				slog.Error("Error writing message to the websocket", slog.Any("error", wErr))
			}
		}

	loop:
		for {
			select {
			case <-ctx.Done():
				err = c.CloseNow()
				if err != nil {
					slog.Error("Can't close websocket connection")
				}

				break loop
			default:
				_, data, e := c.Read(ctx)

				if e != nil {
					slog.Error("Error reading from server", slog.Any("error", e))
					continue loop
				}

				m, ircErr := irc.ParseMessage(string(data))
				if ircErr != nil {
					slog.Error("Error parsing message from server", slog.Any("error", ircErr))
					continue loop
				}

				switch m.Command {
				case "PRIVMSG":
					if len(m.Params) != 2 {
						continue
					}
					message := m.Params[1]
					bot.ch <- TwitchEnvelope{
						Sender:  m.Prefix.Name,
						Message: message,
					}
				case "CAP":
					wErr := c.Write(ctx, websocket.MessageText, []byte("JOIN #"+bot.handle))
					if wErr != nil {
						slog.Error("Error joining channel", slog.Any("error", wErr))
					}
				}
			}
		}
	}()

	return nil
}
