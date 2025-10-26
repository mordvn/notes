package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID    string
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	if s.conn == nil {
		var err error
		s.conn, err = sql.Open("sqlite3", "./notes.db")
		if err != nil {
			return err
		}
	}

	cmd := "Create table if not exists notes (id TEXT PRIMARY KEY, title TEXT, body TEXT);"

	ctx, calcel := context.WithTimeout(context.Background(), time.Second*5)
	defer calcel()
	if _, err := s.conn.ExecContext(ctx, cmd); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAllNotes() ([]Note, error) {
	ctx, calcel := context.WithTimeout(context.Background(), time.Second*5)
	defer calcel()

	rows, err := s.conn.QueryContext(ctx, "SELECT id, title, body FROM notes")
	if err != nil {
		return nil, err
	}

	notes := make([]Note, 0)
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Body); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}

	ctx, calcel := context.WithTimeout(context.Background(), time.Second*5)
	defer calcel()

	cmd := "INSERT INTO notes (id, title, body) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET title=excluded.title, body=excluded.body;"

	if _, err := s.conn.ExecContext(ctx, cmd, note.ID, note.Title, note.Body); err != nil {
		return err
	}

	return nil
}
