package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type noteService interface {
	InsertNote(note models.Notes) error
	GetNotesByUserId(id int) ([]models.Notes, error)
	GetNoteByUserIdNoteId(noteId int, userId int) (models.Notes, error)
	DeleteNoteById(id int) error
	UpdateNote(note models.Notes) error
	ValidNoteByUser(noteId int, userId int) bool
}

type NoteHttpHandler struct {
	noteService
	app    config.Config
	helper helpers.Helper
}

func NewNoteHttpHandler(service noteService, app *config.Config, helper helpers.Helper) *NoteHttpHandler {
	return &NoteHttpHandler{
		noteService: service,
		app:         *app,
		helper:      helper,
	}
}

func (h *NoteHttpHandler) Register(router *chi.Mux) {
	router.Get("/api/note/get_all", h.GetAllNotes)
	router.Post("/api/note/new", h.ApiNewSnippet)
	router.Patch("/api/note/update", h.ApiUpdateSnippet)
	router.Delete("/api/note/delete/{id}", h.ApiDeleteSnippet)
}

func (h *NoteHttpHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	var data []models.Notes
	var err error

	if h.app.Session.Exists(r.Context(), "user_id") {
		log.Println(h.app.Session.Get(r.Context(), "user_id"))
		data, err = h.noteService.GetNotesByUserId(h.app.Session.Get(r.Context(), "user_id").(int))
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
func (h *NoteHttpHandler) ApiNewSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	var response Response
	var note models.Notes

	if h.app.Session.Exists(r.Context(), "user_id") {
		log.Println(h.app.Session.Get(r.Context(), "user_id"))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &note)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		}
		if h.helper.ValidateNote(note) {
			response = Response{
				Message: "Note added",
			}
			note.UserID = h.app.Session.Get(r.Context(), "user_id").(int)
			note.RemovalDate = time.Now()
			note.CreatedAt = time.Now()
			note.UpdatedAt = time.Now()
			er := h.noteService.InsertNote(note)
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

func (h *NoteHttpHandler) ApiDeleteSnippet(w http.ResponseWriter, r *http.Request) {

	var response Response

	if h.app.Session.Exists(r.Context(), "user_id") {
		log.Println(h.app.Session.Get(r.Context(), "user_id"))

		idS := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idS)

		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		if h.noteService.ValidNoteByUser(id, h.app.Session.Get(r.Context(), "user_id").(int)) {
			err := h.noteService.DeleteNoteById(id)
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

func (h *NoteHttpHandler) ApiUpdateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting content type application/json", http.StatusUnsupportedMediaType)
		return
	}

	var response Response
	var note models.Notes

	if h.app.Session.Exists(r.Context(), "user_id") {
		log.Println(h.app.Session.Get(r.Context(), "user_id"))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &note)
		note.UserID = h.app.Session.Get(r.Context(), "user_id").(int)
		if err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		}
		if h.noteService.ValidNoteByUser(note.ID, note.UserID) {
			if h.helper.ValidateUpdateNote(note) {
				response = Response{
					Message: "Note updated",
				}
				er := h.noteService.UpdateNote(note)
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
