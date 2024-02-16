package orchrepository

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/agent"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/parser"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
)

type OrchRepository struct {
	Log   *slog.Logger
	Store storage.Storer
	Agent agent.Agenter
}

func New(log *slog.Logger, store storage.Storer, agent agent.Agenter) *OrchRepository {
	return &OrchRepository{
		Log:   log,
		Store: store,
		Agent: agent,
	}
}

func (o *OrchRepository) CheckExpression(expression string) bool {
	// Стэк для проверки скобок и корректности операндов
	stack := make([]rune, 0)

	for _, char := range expression {
		if char == '(' {
			stack = append(stack, '(')
		} else if char == ')' {
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				// некорректно расставлены скобки
				return false
			}
			stack = stack[:len(stack)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			if len(stack) != 0 && (stack[len(stack)-1] == '+' || stack[len(stack)-1] == '-' || stack[len(stack)-1] == '*' || stack[len(stack)-1] == '/') {
				// повторяющийся операнд
				return false
			}
			stack = append(stack, char)
		} else if char >= '0' && char <= '9' {
			continue
		} else {
			return false // неизвестный символ
		}
	}

	if len(stack) != 0 {
		return false // некорректно расставлены скобки
	}

	return true
}

func (o *OrchRepository) CheckExpOnDb(expression string) (bool, error) {

	//Проверка значения в базе. И если есть, то какой его статус

	return true, nil
}

func (o *OrchRepository) AddExpression(expression string) (int, error) {

	//Добавляем выражение в базу
	return o.Store.ExpressionSave(expression)
}

func (o *OrchRepository) Processing(log *slog.Logger, interval time.Duration, done chan struct{}) {

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
			} else {
				// обновляем статус
				err = o.Store.UpdateStatus(exp, storage.StatusProcess)
				if err != nil {
					println("error on DB")
				}
			}

			// парсинг выражения на части
			parsExp, err := parser.Parsing(exp.ID, exp.Expression)
			if err != nil {
				println(err.Error())
			}

			for !o.Agent.CheckWorkers() {
				time.Sleep(5 * time.Second)
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
