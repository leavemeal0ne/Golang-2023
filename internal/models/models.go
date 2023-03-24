package models

import (
	"time"
)

type TemplateData struct {
	IsAuthenticated bool
	ErrorValues     *MissingValues
	Payload         interface{}
}

type MissingValues struct {
	Email  bool
	Passwd bool
}

type Users struct {
	ID            int       `json:"-"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type Notes struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	RemovalDate time.Time `json:"removal_date"`
}

func (n Notes) ParseDate(date time.Time) string {
	return date.String()[:10]
}
