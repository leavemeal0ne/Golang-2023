package main

import (
	"errors"
	"github.com/leavemeal0ne/Golang-2023/internal/database"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var db *database.DbRepo

func getDatabase(repo *database.DbRepo) {
	db = repo
}
func Home(w http.ResponseWriter, r *http.Request) {
	var data []models.Notes
	var err error = errors.New("new User")
	if session.Exists(r.Context(), "user_id") {
		log.Println(session.Get(r.Context(), "user_id"))
		data, err = db.GetNotesByUserId(session.Get(r.Context(), "user_id").(int))
		template_data.IsAuthenticated = true
	}
	if err != nil {
		log.Println(err)
		log.Println(r.RemoteAddr)
	} else {
		log.Println(" data from session ", data)
	}
	template_data.Payload = data
	RenderTemplate(w, "home_page.html", &template_data)
}

func NewSnippet(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "create_snippet_page.html", &template_data)
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	note := models.Notes{
		UserID:      session.Get(r.Context(), "user_id").(int),
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		RemovalDate: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := db.InsertNote(note)
	if err != nil {
		log.Println(err)
		log.Println("insert error")
	}
	RenderTemplate(w, "create_snippet_page.html", &template_data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	template_data.ErrorValues = nil
	RenderTemplate(w, "login_page.html", &template_data)
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	user, err := db.GetUserByEmail(r.FormValue("email"))
	template_data.ErrorValues = &models.MissingValues{Email: true, Passwd: true}
	if err != nil {
		template_data.ErrorValues.Email = false
		RenderTemplate(w, "login_page.html", &template_data)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(r.FormValue("passwd")))
	if err != nil {
		template_data.ErrorValues.Passwd = false
		RenderTemplate(w, "login_page.html", &template_data)
		return
	}

	app.Session.Put(r.Context(), "user_id", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	template_data.IsAuthenticated = false
	log.Println()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
