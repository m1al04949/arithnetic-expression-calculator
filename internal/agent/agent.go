package agent

import (
	"errors"
	"time"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/parser"
)

type Agent struct {
	Workers int
	TimeSum time.Duration
	TimeSub time.Duration
	TimeMul time.Duration
	TimeDiv time.Duration
}

type Results struct {
	ID     int
	Result float64
}

type Agenter interface {
	CalculateExpression(parser.ParsingExpression) chan Results
	Operation(a, b int, oper string) (float64, error)
}

func New(w int, timeSum, timeSub, timeMul, timeDiv time.Duration) Agenter {
	return &Agent{
		Workers: w,
		TimeSum: timeSum,
		TimeSub: timeSub,
		TimeMul: timeMul,
		TimeDiv: timeDiv,
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
				val, err := ag.Operation(int(a), int(b), token.Value.(string))
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

func (ag *Agent) Operation(a, b int, oper string) (float64, error) {
	switch oper {
	case "+":
		time.Sleep(ag.TimeSum)
		return float64(a + b), nil
	case "-":
		time.Sleep(ag.TimeSub)
		return float64(a - b), nil
	case "*":
		time.Sleep(ag.TimeMul)
		return float64(a * b), nil
	case "/":
		time.Sleep(ag.TimeDiv)
		return float64(a / b), nil
	default:
		return 0, errors.New("Unknown operator: " + oper)
	}
}
