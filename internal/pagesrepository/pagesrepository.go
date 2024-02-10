package pagesrepository

import (
	"log/slog"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/templates"
)

type PagesRepository struct {
	Log       *slog.Logger
	Templates *templates.Template
}

func New(log *slog.Logger, templates *templates.Template) *PagesRepository {
	return &PagesRepository{
		Log:       log,
		Templates: templates,
	}
}
