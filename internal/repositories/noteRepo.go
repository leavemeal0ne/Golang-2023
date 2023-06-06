package repositories

import (
	"context"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"log"
	"time"
)

type NoteRepository struct {
	DB *driver.DB
}

func NewNoteRepository(DB *driver.DB) *NoteRepository {
	return &NoteRepository{DB}
}

func (m *NoteRepository) InsertNote(note models.Notes) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into notes (user_id, title, content, removal_date, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.SQL.ExecContext(ctx, stmt,
		note.UserID,
		note.Title,
		note.Content,
		note.RemovalDate,
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *NoteRepository) GetNotesByUserId(id int) ([]models.Notes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var notes []models.Notes

	query := `
				select 
					id, user_id, title, content, removal_date
				from 
				    notes
				where 
					user_id = $1`

	rows, err := m.DB.SQL.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var note models.Notes
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.RemovalDate,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (m *NoteRepository) GetNoteByUserIdNoteId(noteId int, userId int) (models.Notes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var note models.Notes
	note.ID = noteId

	query := `
				select 
					user_id
				from 
				    notes
				where 
					id = $1`

	row := m.DB.SQL.QueryRowContext(ctx, query, noteId)
	err := row.Scan(
		&note.UserID,
	)
	log.Println(note)
	if err != nil {
		log.Println(err)
		return models.Notes{}, err
	}
	if note.UserID == userId {
		return note, nil
	} else {
		return note, nil
	}

}

func (m *NoteRepository) DeleteNoteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from notes where id = $1`

	_, err := m.DB.SQL.ExecContext(ctx, query, id)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func (m *NoteRepository) UpdateNote(note models.Notes) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update notes set content = $1 where id = $2`

	_, err := m.DB.SQL.ExecContext(ctx, query, note.Content, note.ID)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
