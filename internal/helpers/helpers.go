package helpers

import (
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"net/http"
)

type Helper struct {
	app *config.Config
}

func NewHelper(conf *config.Config) Helper {
	return Helper{
		app: conf,
	}
}

func (m *Helper) IsAuthenticated(r *http.Request) bool {
	exists := m.app.Session.Exists(r.Context(), "user_id")
	return exists
}

func (m *Helper) ValidateNote(note models.Notes) bool {
	return len(note.Content) > 0 && len(note.Title) > 0
}

func (m *Helper) ValidateUpdateNote(note models.Notes) bool {
	return len(note.Content) > 0
}
