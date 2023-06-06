package services

import (
	"github.com/jackc/pgx/v4"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
)

type NoteRepository interface {
	InsertNote(note models.Notes) error
	GetNotesByUserId(id int) ([]models.Notes, error)
	GetNoteByUserIdNoteId(noteId int, userId int) (models.Notes, error)
	DeleteNoteById(id int) error
	UpdateNote(note models.Notes) error
}

type NoteService struct {
	noteRepository NoteRepository
}

func NewNoteService(noteRepository NoteRepository) *NoteService {
	return &NoteService{noteRepository}
}

func (s *NoteService) InsertNote(note models.Notes) error {
	err := s.noteRepository.InsertNote(note)
	return err
}

func (s *NoteService) GetNotesByUserId(id int) ([]models.Notes, error) {
	notes, err := s.noteRepository.GetNotesByUserId(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Notes{}, ErrNotFound
		}
		return []models.Notes{}, err
	}
	return notes, nil
}

func (s *NoteService) GetNoteByUserIdNoteId(noteId int, userId int) (models.Notes, error) {
	note, err := s.noteRepository.GetNoteByUserIdNoteId(noteId, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Notes{}, ErrNotFound
		}
		return models.Notes{}, err
	}
	return note, nil
}

func (s *NoteService) DeleteNoteById(id int) error {
	err := s.noteRepository.DeleteNoteById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *NoteService) UpdateNote(note models.Notes) error {
	err := s.noteRepository.UpdateNote(note)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *NoteService) ValidNoteByUser(noteId int, userId int) bool {

	note, err := s.noteRepository.GetNoteByUserIdNoteId(noteId, userId)
	if err != nil {
		return false
	}
	if note.UserID == userId {
		return true
	}
	return false

}
