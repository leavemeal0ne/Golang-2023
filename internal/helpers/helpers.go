package helpers

import (
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/database"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"net/http"
)

var app config.Config
var db *database.DbRepo

func GetDatabase(repo *database.DbRepo) {
	db = repo
}
func HelpersConfig(conf *config.Config) {
	app = *conf
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

func ValidateNote(note models.Notes) bool {
	return len(note.Content) > 0 && len(note.Title) > 0
}

func ValidNoteByUser(note_id int, user_id int) bool {
	note, err := db.GetNoteByUserIdNoteId(note_id, user_id)
	if err != nil {
		return false
	}
	if note.UserID == user_id {
		return true
	}
	return false
}

func ValidateUpdateNote(note models.Notes) bool {
	return len(note.Content) > 0
}
