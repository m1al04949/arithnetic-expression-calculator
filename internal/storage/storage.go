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
	GetAllExpressions() ([]model.ExpressionTab, error)
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
		status TEXT NOT NULL,
		result FLOAT DEFAULT 0);
		`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ExpressionSave(expression string) (int, error) {

	const op = "storage.ExpressionSave"

	m := &model.ExpressionTab{
		Expression: expression,
		Status:     StatusNew,
	}

	stmt, err := s.db.Prepare("INSERT INTO expressions(expression, status) VALUES ($1, $2);")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Expression, m.Status)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = s.db.QueryRow("SELECT lastval()").Scan(&m.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return m.ID, nil
}

func (s *Storage) GetAllExpressions() ([]model.ExpressionTab, error) {

	const op = "storage.getallexpressions"

	rows, err := s.db.Query("SELECT id, added_at, expression, status, result FROM expressions")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Создаем слайс для хранения результатов
	var expressions []model.ExpressionTab

	// Итерируем по результатам запроса
	for rows.Next() {
		var expression model.ExpressionTab
		err := rows.Scan(&expression.ID, &expression.Added, &expression.Expression, &expression.Status, &expression.Result)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		expressions = append(expressions, expression)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expressions, nil
}

func (s *Storage) GetNewExpression() (model.ExpressionTab, error) {

	// Делаем выборку выражение с status = "New"

	return model.ExpressionTab{}, nil
}

func (s *Storage) UpdateStatus(exp model.ExpressionTab, newStatus string) error {

	// Обновляем статус у выражения

	return nil
}
