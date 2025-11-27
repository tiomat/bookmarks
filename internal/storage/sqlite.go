package storage

import (
	"bookmarks/internal/models"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
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
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
	`
	_, err := db.Exec(query)
	return err
}

func (db *DB) CreateBookmark(b *models.Bookmark) (int64, error) {
	query := `INSERT INTO bookmarks (url, title, excerpt, content, comment, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, b.URL, b.Title, b.Excerpt, b.Content, b.Comment, time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) GetBookmark(id int64) (*models.Bookmark, error) {
	query := `SELECT id, url, title, excerpt, content, comment, created_at FROM bookmarks WHERE id = ?`
	row := db.QueryRow(query, id)

	b := &models.Bookmark{}
	err := row.Scan(&b.ID, &b.URL, &b.Title, &b.Excerpt, &b.Content, &b.Comment, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (db *DB) ListBookmarks() ([]models.Bookmark, error) {
	query := `SELECT id, url, title, excerpt, comment, created_at FROM bookmarks ORDER BY created_at DESC`
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

func (db *DB) UpdateComment(id int64, comment string) error {
	query := `UPDATE bookmarks SET comment = ? WHERE id = ?`
	_, err := db.Exec(query, comment, id)
	return err
}
