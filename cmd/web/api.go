package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"time"
)

type Response struct {
	Message interface{} `json:"response"`
}

func SignUpJson(w http.ResponseWriter, r *http.Request) {
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
	} else if db.IsContainsUserByEmail(user.Email) {
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

	err = db.InsertUser(user)
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

func SignIn(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
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

	result, err := db.GetUserByEmailAndPassword(user.Email, user.Password_hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "Wrong data",
		}
	}
	app.Session.Put(r.Context(), "user_id", result.ID)
	log.Println(session.Get(r.Context(), "user_id"))
	response = Response{
		Message: "Auth complete!",
	}
	template_data.IsAuthenticated = true
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

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}
	var data []models.Notes
	var err error

	if session.Exists(r.Context(), "user_id") {
		log.Println(session.Get(r.Context(), "user_id"))
		data, err = db.GetNotesByUserId(session.Get(r.Context(), "user_id").(int))
		template_data.IsAuthenticated = true
	} else {
		response := Response{
			Message: "need authentication firstly",
		}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
			return
		}

	}

	jsonResponse, err := json.Marshal(data)
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

func ApiNewSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	var response Response
	var note models.Notes

	if session.Exists(r.Context(), "user_id") {
		log.Println(session.Get(r.Context(), "user_id"))
		template_data.IsAuthenticated = true

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &note)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		}
		if helpers.ValidateNote(note) {
			response = Response{
				Message: "Note added",
			}
			note.UserID = session.Get(r.Context(), "user_id").(int)
			note.RemovalDate = time.Now()
			note.CreatedAt = time.Now()
			note.UpdatedAt = time.Now()
			er := db.InsertNote(note)
			if er != nil {
				log.Println(er)
				response = Response{
					Message: "DB error",
				}
			}
		} else {
			response = Response{
				Message: "Invalid data",
			}
		}
	} else {
		response = Response{
			Message: "need authentication firstly",
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
}

func ApiDeleteSnippet(w http.ResponseWriter, r *http.Request) {

	var response Response

	if session.Exists(r.Context(), "user_id") {
		log.Println(session.Get(r.Context(), "user_id"))
		template_data.IsAuthenticated = true

		idS := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idS)

		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		if helpers.ValidNoteByUser(id, session.Get(r.Context(), "user_id").(int)) {
			err := db.DeleteNoteById(id)
			if err != nil {
				http.Error(w, "Failed to parse JSON data", http.StatusInternalServerError)
				response = Response{
					Message: "DB error",
				}
			} else {
				response = Response{
					Message: id,
				}
				w.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			response = Response{
				Message: "Invalid id",
			}
		}
	} else {
		response = Response{
			Message: "need authentication firstly",
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
}

//	{
//		"id"
//		"content"
//	}
func ApiUpdateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	var response Response
	var note models.Notes

	if session.Exists(r.Context(), "user_id") {
		log.Println(session.Get(r.Context(), "user_id"))
		template_data.IsAuthenticated = true

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &note)
		note.UserID = session.Get(r.Context(), "user_id").(int)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		}
		if helpers.ValidNoteByUser(note.ID, note.UserID) {
			if helpers.ValidateUpdateNote(note) {
				response = Response{
					Message: "Note updated",
				}
				er := db.UpdateNote(note)
				if er != nil {
					log.Println(er)
					response = Response{
						Message: "DB error",
					}
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				response = Response{
					Message: "Note updated",
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response = Response{
				Message: "Invalid data",
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{
			Message: "need authentication firstly",
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
}
