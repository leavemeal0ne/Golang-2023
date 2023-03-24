package database

import (
	"context"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type DbRepo struct {
	DB *driver.DB
}

func (m *DbRepo) InsertUser(user models.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password_hash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `insert into users(email, password_hash,created_at, updated_at) values ($1, $2, $3, $4)`

	_, err = m.DB.SQL.ExecContext(ctx, stmt,
		user.Email,
		hash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *DbRepo) InsertNote(note models.Notes) error {
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

func (m *DbRepo) GetNotesByUserId(id int) ([]models.Notes, error) {
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

func (m *DbRepo) GetUserByEmail(email string) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.Users
	user.Email = email

	query := `
				select 
					id, password_hash
				from 
				    users
				where 
					email = $1`

	row := m.DB.SQL.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Password_hash,
	)
	if err != nil {
		log.Println(err)
		return models.Users{}, err
	}

	return user, nil
}

func (m *DbRepo) IsContainsUserByEmail(email string) bool {
	_, err := m.GetUserByEmail(email)
	if err != nil {
		return false
	}
	return true
}

func (m *DbRepo) GetUserByEmailAndPassword(email, password string) (models.Users, error) {
	user, err := m.GetUserByEmail(email)
	if err != nil {
		return models.Users{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(password))
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (m *DbRepo) GetNoteByUserIdNoteId(note_id int, user_id int) (models.Notes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var note models.Notes
	note.ID = note_id

	query := `
				select 
					user_id
				from 
				    notes
				where 
					id = $1`

	row := m.DB.SQL.QueryRowContext(ctx, query, note_id)
	err := row.Scan(
		&note.UserID,
	)
	log.Println(note)
	if err != nil {
		log.Println(err)
		return models.Notes{}, err
	}
	if note.UserID == user_id {
		return note, nil
	} else {
		return note, nil
	}

}

func (m *DbRepo) DeleteNoteById(id int) error {
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

func (m *DbRepo) UpdateNote(note models.Notes) error {
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
