package usersrepository

import (
	"log/slog"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
	"github.com/m1al04949/arithnetic-expression-calculator/pkg/hash"
)

type UsersRepository struct {
	Log    *slog.Logger
	Config *config.Config
	Store  storage.Storer
}

func New(log *slog.Logger, cfg *config.Config, store storage.Storer) *UsersRepository {
	return &UsersRepository{
		Log:    log,
		Config: cfg,
		Store:  store,
	}
}

func (u *UsersRepository) CheckUserOnDb(user string) (bool, error) {

	exists, err := u.Store.CheckUserExists(user)
	if err != nil {
		u.Log.Error("server internal error", slog.Any("error", err.Error()))
	}

	return exists, err
}

func (u *UsersRepository) CreateUser(user, password string) error {

	err := u.Store.CreateUser(user, password)
	if err != nil {
		u.Log.Error("server internal error", slog.Any("error", err.Error()))
	}

	return nil
}

func (u *UsersRepository) CheckAuthorization(user, password string) (bool, error) {

	check, pass, err := u.Store.CheckAuth(user)
	if err != nil {
		u.Log.Error("server internal error", slog.Any("error", err.Error()))
	}

	check = hash.ComparePasswordWithHash(password, pass)

	return check, nil
}
