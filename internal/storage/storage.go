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
	UpdateResult(int, float64) error
	RefreshStatus() error
	CheckUserExists(string) (bool, error)
	CreateUser(string, string) error
	CheckAuth(string) (bool, int, string, error)
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
		user_name TEXT NOT NULL,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expression TEXT,
		status TEXT NOT NULL,
		result FLOAT DEFAULT 0);
		`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.RefreshStatus()
	if err != nil {
		println(err.Error())
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		user_id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
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

	// const op = "storage.getnewexpression"

	// Делаем выборку выражение с status = "new"
	var row model.ExpressionTab

	err := s.db.QueryRow("SELECT * FROM expressions WHERE status = $1 ORDER BY id LIMIT 1", StatusNew).Scan(&row.ID, &row.Added, &row.Expression, &row.Status, &row.Result)
	if err != nil {
		return model.ExpressionTab{}, err
	}

	return row, nil
}

func (s *Storage) UpdateStatus(exp model.ExpressionTab, newStatus string) error {

	const op = "storage.updatestatus"

	// Обновляем статус у выражения
	_, err := s.db.Exec("UPDATE expressions SET status = $1 WHERE id = $2", newStatus, exp.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateResult(ID int, result float64) error {

	const op = "storage.updateresult"

	// Обновляем статус у выражения
	_, err := s.db.Exec("UPDATE expressions SET result = $2, status = $3 WHERE id = $1", ID, result, StatusComplete)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RefreshStatus() error {

	const op = "storage.refreshstatus"

	rows, err := s.db.Query("SELECT id, added_at, expression, status, result FROM expressions WHERE status = $1", StatusProcess)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Итерируем по результатам запроса
	for rows.Next() {
		var exp model.ExpressionTab
		err := rows.Scan(&exp.ID, &exp.Added, &exp.Expression, &exp.Status, &exp.Result)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = s.UpdateStatus(exp, StatusNew)
		if err != nil {
			println(err.Error())
			return fmt.Errorf("%s: %w", op, err)
		}
		err = rows.Err()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) CheckUserExists(username string) (bool, error) {

	const op = "storage.checkuserexists"

	var count int
	// Проверяем наличие пользователя
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1", username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (s *Storage) CreateUser(user, password string) error {

	const op = "storage.createuser"

	// Добавляем пользователя
	_, err := s.db.Exec(`
        INSERT INTO users (login, password_hash)
        VALUES ($1, $2)
    `, user, password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CheckAuth(user string) (bool, int, string, error) {

	const op = "storage.checkauth"

	var (
		id       int
		passOnDB string
	)

	// Проверяем пользователя
	err := s.db.QueryRow(`
        SELECT user_id, password_hash FROM users WHERE login = $1
    `, user).Scan(&id, &passOnDB)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, "", nil
		}
		return false, 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return true, id, passOnDB, nil
}
