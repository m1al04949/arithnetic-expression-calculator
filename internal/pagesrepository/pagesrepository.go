package pagesrepository

import (
	"log/slog"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/agent"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/templates"
)

type PagesRepository struct {
	Log       *slog.Logger
	Templates *templates.Template
	Config    *config.Config
	Store     storage.Storer
	Agent     agent.Agenter
}

func New(log *slog.Logger, templates *templates.Template, cfg *config.Config, store storage.Storer, agent agent.Agenter) *PagesRepository {
	return &PagesRepository{
		Log:       log,
		Templates: templates,
		Config:    cfg,
		Store:     store,
		Agent:     agent,
	}
}
