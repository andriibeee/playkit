package main

import (
	"context"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"playlist/auth"
	"playlist/playlist"
	"playlist/shared"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vk-rv/pvx"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	cfg := shared.NewConfiguration()

	h := playlist.NewHub()

	if strings.TrimSpace(cfg.DB) == "" {
		slog.Error("DB configuration is empty")
		return
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DB)
	if err != nil {
		slog.Error("Unable to create connection pool", slog.Any("error", err))
		return
	}

	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		slog.Error("Unable to connect to database", slog.Any("error", err))
		return
	}

	k, err := hex.DecodeString(cfg.AuthSigningKey)
	if err != nil {
		slog.Error("Unable to decode hex string", slog.Any("error", err))
		return
	}

	symK := pvx.NewSymmetricKey(k, pvx.Version4)
	pv4 := pvx.NewPV4Local()

	ps := playlist.NewPlaylistService(dbpool)
	api := playlist.NewYoutubeAPI(cfg.YoutubeAPIKey)
	uc := playlist.NewPlaylistUseCase(h, ps, api)
	s := playlist.NewStreamer(h, uc, symK, pv4)

	amw := shared.NewAuthMiddleware(pv4, symK)

	p := playlist.NewPlaylistPort(uc, s, amw)
	ap := auth.NewAuthPort(cfg.OAuth2ID, cfg.OAuth2Secret, cfg.OAuth2RedirectURL, pv4, symK)

	r := chi.NewRouter()

	r.Route("/auth", ap.Router)
	r.Route("/playlist", p.Router)

	port := cfg.Port
	if port == "" {
		port = ":3000"
	} else if port[:0] != ":" {
		port = ":" + cfg.Port
	}

	srv := http.Server{
		Addr:    port,
		Handler: r,

		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			slog.Error("Unable to start server", slog.Any("error", err))
		}
	}()

	slog.Info("Listening", slog.String("port", port))

	<-c

	err = srv.Shutdown(context.Background())
	if err != nil {
		slog.Error("Unable to shutdown server", slog.Any("error", err))
	}
}
