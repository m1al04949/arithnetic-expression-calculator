package templates

import (
	"html/template"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
)

type Template struct {
	Main template.Template
}

func New(cfg config.Config) (*Template, error) {

	mainPage, err := template.ParseFiles("путь до шаблона")
	if err != nil {
		return nil, err
	}

	return &Template{
		Main: *mainPage,
	}, nil
}
