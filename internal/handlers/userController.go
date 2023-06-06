package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"io"
	"log"
	"net/http"
	"net/mail"
)

type userService interface {
	InsertUser(user models.Users) error
	GetUserByEmail(email string) (models.Users, error)
	IsContainsUserByEmail(email string) bool
	GetUserByEmailAndPassword(email, password string) (models.Users, error)
}

type UserHttpHandler struct {
	userService
	app config.Config
}

func NewUserHttpHandler(service userService, app *config.Config) *UserHttpHandler {
	return &UserHttpHandler{
		userService: service,
		app:         *app,
	}
}

func (h *UserHttpHandler) Register(router *chi.Mux) {
	router.Post("/api/user/signup", h.SignUpJson)
	router.Post("/api/user/login", h.SignIn)
}

type Response struct {
	Message interface{} `json:"response"`
}

func (h *UserHttpHandler) SignUpJson(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var user models.Users
	err = json.Unmarshal(body, &user)

	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(user.Email)

	var response Response

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "Mail should be valid",
		}
	} else if len(user.Password_hash) <= 5 {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "too short password",
		}
	} else if h.userService.IsContainsUserByEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "Email is used",
		}
	} else {
		w.WriteHeader(http.StatusOK)
		response = Response{
			Message: "User created",
		}
	}

	err = h.userService.InsertUser(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	jsonResponse, err := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Send the response
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

}

func (h *UserHttpHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	_ = h.app.Session.RenewToken(r.Context())
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var user models.Users
	err = json.Unmarshal(body, &user)

	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(user.Email)

	var response Response

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "Mail should be valid",
		}
	}

	result, err := h.userService.GetUserByEmailAndPassword(user.Email, user.Password_hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "Wrong data",
		}
	}
	h.app.Session.Put(r.Context(), "user_id", result.ID)
	log.Println(h.app.Session.Get(r.Context(), "user_id"))
	response = Response{
		Message: "Auth complete!",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Send the response
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
}
