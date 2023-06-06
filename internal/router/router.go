package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/middleware"
	"net/http"
)

func Routes() *chi.Mux {

	mux := chi.NewRouter()
	mux.Use(middleware.SessionLoad)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
