package storage

import (
	"bookmarks/internal/models"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT,
		excerpt TEXT,
		content TEXT,
		comment TEXT,
		deleted BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Migration: Try to add 'deleted' column if it doesn't exist (for existing DBs)
	// We ignore the error because if the column exists, it will fail, which is fine.
	migration := `ALTER TABLE bookmarks ADD COLUMN deleted BOOLEAN DEFAULT 0;`
	db.Exec(migration)

	return nil
}

func (db *DB) CreateBookmark(b *models.Bookmark) (int64, error) {
	query := `INSERT INTO bookmarks (url, title, excerpt, content, comment, deleted, created_at) VALUES (?, ?, ?, ?, ?, 0, ?)`
	res, err := db.Exec(query, b.URL, b.Title, b.Excerpt, b.Content, b.Comment, time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) GetBookmark(id int64) (*models.Bookmark, error) {
	query := `SELECT id, url, title, excerpt, content, comment, deleted, created_at FROM bookmarks WHERE id = ?`
	row := db.QueryRow(query, id)

	b := &models.Bookmark{}
	err := row.Scan(&b.ID, &b.URL, &b.Title, &b.Excerpt, &b.Content, &b.Comment, &b.Deleted, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (db *DB) ListBookmarks() ([]models.Bookmark, error) {
	query := `SELECT id, url, title, excerpt, comment, created_at FROM bookmarks WHERE deleted = 0 ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		// Note: We are not fetching content here to keep the list lightweight
		err := rows.Scan(&b.ID, &b.URL, &b.Title, &b.Excerpt, &b.Comment, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func (db *DB) DeleteBookmark(id int64) error {
	query := `UPDATE bookmarks SET deleted = 1 WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) UpdateComment(id int64, comment string) error {
	query := `UPDATE bookmarks SET comment = ? WHERE id = ?`
	_, err := db.Exec(query, comment, id)
	return err
}
