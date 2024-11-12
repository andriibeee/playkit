package playlist

import (
	"encoding/json"
	"net/http"
	"playlist/shared"

	"github.com/go-chi/chi/v5"
)

type PlaylistPort struct {
	uc  *PlaylistUseCase
	s   *Streamer
	amw *shared.AuthMiddleware
}

func NewPlaylistPort(uc *PlaylistUseCase, s *Streamer, amw *shared.AuthMiddleware) *PlaylistPort {
	return &PlaylistPort{
		uc:  uc,
		s:   s,
		amw: amw,
	}
}

func (p *PlaylistPort) playlist(w http.ResponseWriter, r *http.Request) {
	handle := shared.ExtractHandle(r.Context())

	pl, err := p.uc.GetPlaylist(r.Context(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(pl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlaylistPort) deleteVideoFromPlaylist(w http.ResponseWriter, r *http.Request) {
	handle := shared.ExtractHandle(r.Context())

	err := p.uc.DeleteVideo(r.Context(), handle, chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PlaylistPort) skipVideo(w http.ResponseWriter, r *http.Request) {
	handle := shared.ExtractHandle(r.Context())

	err := p.uc.SkipVideo(r.Context(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PlaylistPort) Router(r chi.Router) {
	r.Handle("/stream", p.s)
	r.Group(func(r chi.Router) {
		r.Use(p.amw.Middleware)
		r.Get("/", p.playlist)
		r.Post("/skip", p.skipVideo)
		r.Delete("/{id}", p.deleteVideoFromPlaylist)
	})
}
