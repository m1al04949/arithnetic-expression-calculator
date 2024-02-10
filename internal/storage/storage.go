package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/model"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

type Storer interface {
	Close()
	CreateTabs() error
	Open() error
	ExpressionSave(string) (int, error)
	GetNewExpression() (model.ExpressionTab, error)
	UpdateStatus(model.ExpressionTab, string) error
}

var (
	ErreExpExists = errors.New("expression exists")
)

const (
	StatusNew      = "new"
	StatusProcess  = "in processing"
	StatusComplete = "completed"
)

// Get instance
func New(cfgpath, dburl string) Storer {
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
		expression TEXT,
		status TEXT NOT NULL);
		`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ExpressionSave(expression string) (int, error) {

	// const op = "storage.SaveUser"

	// m := &model.Users{
	// 	UserID: userToSave,
	// }

	// if err := s.db.QueryRow("SELECT (created_at) FROM users WHERE user_id=$1",
	// 	userToSave).Scan(&m.CreatedAt); err != nil {
	// 	stmt, err := s.db.Prepare("INSERT INTO users(user_id) VALUES ($1)")
	// 	if err != nil {
	// 		return fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	defer stmt.Close()

	// 	_, err = stmt.Exec(userToSave)
	// 	if err != nil {
	// 		if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
	// 			return fmt.Errorf("%s: %w, created at %s", op, ErrUserExists, m.CreatedAt)
	// 		}
	// 		return fmt.Errorf("%s: %w", op, err)
	// 	}
	// } else {
	// 	return fmt.Errorf("%s: %w, created at %s", op, ErrUserExists, m.CreatedAt)
	// }

	return 0, nil
}

func (s *Storage) GetNewExpression() (model.ExpressionTab, error) {

	// Делаем выборку выражение с status = "New"

	return model.ExpressionTab{}, nil
}

func (s *Storage) UpdateStatus(exp model.ExpressionTab, newStatus string) error {

	// Обновляем статус у выражения

	return nil
}
