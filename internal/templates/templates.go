package templates

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
)

type Template struct {
	Auth        template.Template
	Main        template.Template
	Register    template.Template
	Login       template.Template
	Settings    template.Template
	Expressions template.Template
	Tasks       template.Template
}

func New(cfg config.Config) (*Template, error) {
	// Path to HTML-file
	authPath := "/internal/templates/main/auth.html"
	mainPath := "/internal/templates/main/index.html"
	settingPath := "/internal/templates/settings/settings.html"
	expressionsPath := "/internal/templates/expressions/expressions.html"
	tasksPath := "/internal/templates/tasks/tasks.html"
	regPath := "/internal/templates/register/register.html"
	loginPath := "/internal/templates/login/login.html"

	// Get Template Path
	currDir, _ := os.Getwd()
	projectDir := filepath.Join(currDir, "..", "..")

	templateAuthPath := filepath.Join(projectDir, authPath)
	authPage, err := template.ParseFiles(templateAuthPath)
	if err != nil {
		return nil, err
	}

	templateMainPath := filepath.Join(projectDir, mainPath)
	mainPage, err := template.ParseFiles(templateMainPath)
	if err != nil {
		return nil, err
	}

	templateRegPath := filepath.Join(projectDir, regPath)
	regPage, err := template.ParseFiles(templateRegPath)
	if err != nil {
		return nil, err
	}

	templateLoginPath := filepath.Join(projectDir, loginPath)
	loginPage, err := template.ParseFiles(templateLoginPath)
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

	templateTasksPath := filepath.Join(projectDir, tasksPath)
	tasksPage, err := template.ParseFiles(templateTasksPath)
	if err != nil {
		return nil, err
	}

	return &Template{
		Auth:        *authPage,
		Main:        *mainPage,
		Register:    *regPage,
		Login:       *loginPage,
		Settings:    *settingsPage,
		Expressions: *expressionsPage,
		Tasks:       *tasksPage,
	}, nil
}
