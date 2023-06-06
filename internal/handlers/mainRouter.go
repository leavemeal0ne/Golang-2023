package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UserHttpHandlerEndpoints interface {
	SignUpJson(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Register(router *chi.Mux)
}

type NoteHttpHandlerEndpoints interface {
	GetAllNotes(w http.ResponseWriter, r *http.Request)
	ApiNewSnippet(w http.ResponseWriter, r *http.Request)
	ApiDeleteSnippet(w http.ResponseWriter, r *http.Request)
	ApiUpdateSnippet(w http.ResponseWriter, r *http.Request)
	Register(router *chi.Mux)
}

func SetupRoutes(router *chi.Mux, userHttp UserHttpHandlerEndpoints, noteHttp NoteHttpHandlerEndpoints) {
	userHttp.Register(router)
	noteHttp.Register(router)
}
