package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"net/http"
)

func routes(app *config.Config) http.Handler {

	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", Home)
	mux.Get("/user/login", Login)
	mux.Get("/user/logout", Logout)
	mux.Post("/user/login", Authentication)

	//api
	mux.Get("/api/user/signup", SignUpJson)
	mux.Get("/api/user/login", SignIn)
	mux.Get("/api/snippet/get_all", GetAllNotes)
	mux.Post("/api/snippet/new", ApiNewSnippet)
	mux.Patch("/api/snippet/update", ApiUpdateSnippet)
	mux.Delete("/api/snippet/delete/{id}", ApiDeleteSnippet)

	mux.Route("/snippet", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/new", NewSnippet)
		r.Post("/new", CreateSnippet)
	})
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
