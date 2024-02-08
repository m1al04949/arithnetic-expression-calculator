package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

var (
	ErreExpExists = errors.New("expression exists")
)

// Get instance
func New(cfgpath, dburl string) *Storage {
	return &Storage{
		config: &Config{
			ConfigPath:  cfgpath,
			DatabaseURL: dburl,
		},
	}
}

// Open connection to DB
func (s *Storage) Open() error {

	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Close connection
func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CreateTabs() error {
	const op = "storage.CreateTabs"

	_, err := s.db.Exec(`
	    CREATE TABLE IF NOT EXISTS expressions(
		id SERIAL PRIMARY KEY,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expression TEXT);
		`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
