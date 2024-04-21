package orchrepository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"unicode"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/agent"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/parser"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
)

type OrchRepository struct {
	Cfg   *config.Config
	Log   *slog.Logger
	Store storage.Storer
	Agent agent.Agenter
}

func New(cfg *config.Config, log *slog.Logger, store storage.Storer, agent agent.Agenter) *OrchRepository {
	return &OrchRepository{
		Cfg:   cfg,
		Log:   log,
		Store: store,
		Agent: agent,
	}
}

func (o *OrchRepository) CheckExpression(expression string) (error, bool) {
	prevCharIsOperator := true // Переменная для отслеживания, является ли предыдущий символ оператором

	for _, char := range expression {
		if char == '(' || char == ')' {
			return fmt.Errorf("в выражении содержатся скобки"), false
		} else if !unicode.IsDigit(char) && !unicode.IsSpace(char) {
			if prevCharIsOperator {
				return fmt.Errorf("выражение содержит подряд идущие операторы"), false
			}
			prevCharIsOperator = true
		} else if unicode.IsDigit(char) {
			prevCharIsOperator = false
		}
	}

	return nil, true
}

func (o *OrchRepository) CheckExpOnDb(user, expression string) (bool, error) {

	//Проверка значения в базе. И если есть, то какой его статус

	return true, nil
}

func (o *OrchRepository) AddExpression(user, expression string) (int, error) {

	//Добавляем выражение в базу
	return o.Store.ExpressionSave(user, expression)
}

func (o *OrchRepository) Processing(user string, log *slog.Logger, interval time.Duration, done chan struct{}) {

	// Опрашиваем базу данных на предмет получения новых выражений
	// Парсим новые выражения на части
	// Отправляем задания агенту

	log.Info("processing started")
	resultChan := make(chan agent.Results)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// выполнение запроса к базе данных
			log.Info("check another expression on DB")
			exp, err := o.Store.GetNewExpression()
			if err != nil {
				if err != sql.ErrNoRows {
					return
				}
			}

			// парсинг выражения на части
			parsExp, err := parser.Parsing(exp.ID, exp.Expression)
			if err != nil {
				println(err.Error())
				return
			}

			for {
				_, ok := o.Agent.CheckWorkers()
				if ok {
					break
				}
				time.Sleep(5 * time.Second)
			}

			// обновляем статус
			err = o.Store.UpdateStatus(exp, storage.StatusProcess)
			if err != nil {
				println("error on DB")
				return
			}

			go func() {
				o.Agent.DecrementWorkers()
				resultExp := o.Agent.CalculateExpression(parsExp)
				resultChan <- <-resultExp
			}()

		case result := <-resultChan:
			o.Agent.IncrementWorkers()
			err := o.Store.UpdateResult(result.ID, result.Result)
			if err != nil {
				println(err.Error())
				return
			}
		// Ожидаем завершения работы, если получили сигнал через канал done
		case <-done:
			return
		}
	}
}
