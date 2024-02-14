package parser

import shuntingYard "github.com/mgenware/go-shunting-yard"

type ParsingExpression struct {
	ID          int
	Expressions []*shuntingYard.RPNToken
}

func Parsing(id int, exp string) (ParsingExpression, error) {

	// Обратная польская волына
	infixTokens, err := shuntingYard.Scan(exp)
	if err != nil {
		return ParsingExpression{}, err
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return ParsingExpression{}, err
	}

	return ParsingExpression{
		ID:          id,
		Expressions: postfixTokens,
	}, nil
}
