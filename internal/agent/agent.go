package agent

import (
	"errors"
	"time"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/parser"
)

type Agent struct {
	Config *config.Config
}

type Results struct {
	ID     int
	Result float64
}

type Agenter interface {
	CalculateExpression(parser.ParsingExpression) chan Results
	Operation(a, b float64, oper string) (float64, error)
	DecrementWorkers()
	IncrementWorkers()
	CheckWorkers() (int, bool)
}

func New(cfg *config.Config) Agenter {
	return &Agent{
		Config: cfg,
	}
}

func (ag *Agent) CalculateExpression(parsExp parser.ParsingExpression) chan Results {
	var stack []float64
	out := make(chan Results)

	go func() {
		defer close(out)
		for _, token := range parsExp.Expressions {
			// push all operands to the stack
			if token.Type == 1 {
				val := token.Value.(int)
				stack = append(stack, float64(val))
			} else {
				if len(stack) < 2 {
					return
				}
				a, b := stack[len(stack)-2], stack[len(stack)-1]
				stack = stack[:len(stack)-2]
				//Передача 2 чисел и оператора в вычислитель
				val, err := ag.Operation(a, b, token.Value.(string))
				if err != nil {
					println(err.Error())
					return
				}
				// push result back to stack
				stack = append(stack, val)
			}
		}
		if len(stack) != 1 {
			return
		}

		out <- Results{
			ID:     parsExp.ID,
			Result: stack[0],
		}
	}()

	return out
}

func (ag *Agent) Operation(a, b float64, oper string) (float64, error) {
	switch oper {
	case "+":
		time.Sleep(ag.Config.OperationSumInterval)
		return a + b, nil
	case "-":
		time.Sleep(ag.Config.OperationSubInterval)
		return a - b, nil
	case "*":
		time.Sleep(ag.Config.OperationMulInterval)
		return a * b, nil
	case "/":
		time.Sleep(ag.Config.OperationDivInterval)
		return a / b, nil
	default:
		return 0, errors.New("Unknown operator: " + oper)
	}
}

func (ag *Agent) DecrementWorkers() {
	ag.Config.Workers = ag.Config.Workers - 1
}

func (ag *Agent) IncrementWorkers() {
	ag.Config.Workers = ag.Config.Workers + 1
}

func (ag *Agent) CheckWorkers() (int, bool) {
	return ag.Config.Workers, ag.Config.Workers > 0
}
