package orchrepository

import (
	"log/slog"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/templates"
)

type OrchRepository struct {
	Log       *slog.Logger
	Store     storage.Storer
	Templates *templates.Template
}

func New(log *slog.Logger, store storage.Storer) *OrchRepository {
	return &OrchRepository{
		Log:   log,
		Store: store,
	}
}

func (o *OrchRepository) CheckExpression(expression string) (bool, error) {

	//Обратная польская нотация

	return false, nil
}

func (o *OrchRepository) CheckExpOnDb(expression string) (bool, error) {

	//Проверка значения в базе. И если есть, то какой его статус

	return false, nil
}

func (o *OrchRepository) AddExpression(expression string) (int, error) {

	//Добавляем выражение в базу
	return o.Store.ExpressionSave(expression)
}
