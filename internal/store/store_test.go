package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	s := new(Store)
	require.NoError(t, s.Init())
	require.NotNil(t, s.conn)
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	s := &Store{conn: db}
	require.NoError(t, s.Init())
	return s
}

func TestSaveAndGetAllNotes(t *testing.T) {
	s := newTestStore(t)

	n := Note{Title: "Hello", Body: "World"}
	require.NoError(t, s.SaveNote(n))

	notes, err := s.GetAllNotes()
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, n.Title, notes[0].Title)
	require.Equal(t, n.Body, notes[0].Body)
}

func TestUpdateNote(t *testing.T) {
	s := newTestStore(t)

	n := Note{ID: "123", Title: "Old", Body: "B"}
	require.NoError(t, s.SaveNote(n))

	n.Title = "New"
	require.NoError(t, s.SaveNote(n))

	notes, err := s.GetAllNotes()
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "New", notes[0].Title)
}
