package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func New(dbPath string) (*Database, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{DB: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS music_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		query TEXT NOT NULL,
		title TEXT NOT NULL,
		file_id TEXT NOT NULL,
		duration INTEGER,
		performer TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(query, title)
	);

	CREATE INDEX IF NOT EXISTS idx_music_cache_query ON music_cache(query);
	CREATE INDEX IF NOT EXISTS idx_music_cache_file_id ON music_cache(file_id);
	`

	_, err := d.DB.Exec(schema)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
		return err
	}

	return nil
}

// CacheMusic stores downloaded music file_id for future use
func (d *Database) CacheMusic(query, title, fileID string, duration int, performer string) error {
	_, err := d.DB.Exec(`
		INSERT OR REPLACE INTO music_cache (query, title, file_id, duration, performer)
		VALUES (?, ?, ?, ?, ?)
	`, query, title, fileID, duration, performer)
	return err
}

// GetCachedMusic retrieves cached music by query and title
func (d *Database) GetCachedMusic(query, title string) (string, error) {
	var fileID string
	err := d.DB.QueryRow(`
		SELECT file_id FROM music_cache
		WHERE query = ? AND title = ?
	`, query, title).Scan(&fileID)
	return fileID, err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
