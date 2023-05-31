package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/database"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
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

func DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	ids := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(ids)
	if helpers.ValidNoteByUser(id, session.Get(r.Context(), "user_id").(int)) {
		err := db.DeleteNoteById(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Authentication(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	log.Println("xd")
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	template_data.ErrorValues = nil
	RenderTemplate(w, "register_page.html", &template_data)
}
func SignUpValidate(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	template_data.ErrorValues = &models.MissingValues{Email: true, Passwd: true}
	_, err := mail.ParseAddress(r.FormValue("email"))
	if err != nil {
		template_data.ErrorValues.Email = false
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "register_page.html", &template_data)
		return
	}
	_, err = db.GetUserByEmail(r.FormValue("email"))
	if err != nil {
		var user models.Users
		user.Email = r.FormValue("email")
		user.Password_hash = r.FormValue("email")
		user.Password_hash = r.FormValue("passwd")
		err := db.InsertUser(user)
		w.WriteHeader(http.StatusOK)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		RenderTemplate(w, "home_page.html", &template_data)
		return
	} else {
		template_data.ErrorValues.Email = false
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "register_page.html", &template_data)
		return
	}
}
