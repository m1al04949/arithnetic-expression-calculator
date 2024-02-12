package templates

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
)

type Template struct {
	Main        template.Template
	Settings    template.Template
	Expressions template.Template
}

func New(cfg config.Config) (*Template, error) {
	// Path to HTML-file
	mainPath := "/internal/templates/main/index.html"
	settingPath := "/internal/templates/settings/settings.html"
	expressionsPath := "/internal/templates/expressions/expressions.html"

	// Get Template Path
	currDir, _ := os.Getwd()
	projectDir := filepath.Join(currDir, "..", "..")

	templateMainPath := filepath.Join(projectDir, mainPath)
	mainPage, err := template.ParseFiles(templateMainPath)
	if err != nil {
		return nil, err
	}

	templateSettingsPath := filepath.Join(projectDir, settingPath)
	settingsPage, err := template.ParseFiles(templateSettingsPath)
	if err != nil {
		return nil, err
	}

	templateExpressionsPath := filepath.Join(projectDir, expressionsPath)
	expressionsPage, err := template.ParseFiles(templateExpressionsPath)
	if err != nil {
		return nil, err
	}

	return &Template{
		Main:        *mainPage,
		Settings:    *settingsPage,
		Expressions: *expressionsPage,
	}, nil
}
