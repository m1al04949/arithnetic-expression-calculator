package orchrepository

import (
	"log/slog"
	"time"

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

func (o *OrchRepository) Processing(interval time.Duration, done chan struct{}) {

	// Опрашиваем базу данных на предмет получения новых выражений
	// Парсим новые выражения на части
	// Отправляем задания агенту

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// выполнение запроса к базе данных
			exp, err := o.Store.GetNewExpression()
			if err != nil {
				println("error on DB")
			}
			// обновляем статус
			err = o.Store.UpdateStatus(exp, storage.StatusProcess)
			if err != nil {
				println("error on DB")
			}

			// парсинг выражения на части
			// expPars := parser.Parsing(exp.Expression)

		// Ожидаем завершения работы, если получили сигнал через канал done
		case <-done:
			return
		}
	}

}
