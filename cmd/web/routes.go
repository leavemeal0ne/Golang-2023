package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", Home)
	mux.Get("/user/login", Login)
	mux.Get("/user/logout", Logout)
	mux.Post("/user/login", Authentication)
	mux.Get("/user/signup", SignUp)
	mux.Post("/user/signup", SignUpValidate)

	//api
	mux.Post("/api/user/signup", SignUpJson)
	mux.Post("/api/user/login", SignIn)
	mux.Get("/api/snippet/get_all", GetAllNotes)
	mux.Post("/api/snippet/new", ApiNewSnippet)
	mux.Patch("/api/snippet/update", ApiUpdateSnippet)
	mux.Delete("/api/snippet/delete/{id}", ApiDeleteSnippet)

	mux.Route("/snippet", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/new", NewSnippet)
		r.Post("/new", CreateSnippet)
		r.Get("/delete/{id}", DeleteSnippet)
	})
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
