package usersrepository

import (
	"log/slog"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
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
