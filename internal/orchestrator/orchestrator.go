package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/agent"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/http-server/handlers/expressions"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/http-server/handlers/pages"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/http-server/handlers/users"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/orchrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/pagesrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/usersrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/templates"
)

type Orchestrator struct {
	Config *config.Config
	Log    *slog.Logger
	Server *http.Server
}

func New(cfg *config.Config, log *slog.Logger) *Orchestrator {
	return &Orchestrator{
		Config: cfg,
		Log:    log,
		Server: &http.Server{},
	}
}

func (o *Orchestrator) RunServer() error {

	// Storage Initializing
	store := storage.New(o.Config.StoragePath, o.Config.DatabaseURL)
	if err := store.Open(); err != nil {
		o.Log.Error("failed to init storage")
		return err
	}
	defer store.Close()
	if err := store.CreateTabs(); err != nil {
		o.Log.Error("failed to init tabs")
		return err
	}
	o.Log.Info("storage is initialized")

	// Templates Initializing
	templates, err := templates.New(*o.Config)
	if err != nil {
		o.Log.Error("failed to init templates")
		return err
	}
	o.Log.Info("templates is initialized")

	// Init Agents
	agent := agent.New(o.Config)
	o.Log.Info("agent is initialized")

	//Init Repositories
	orchRepository := orchrepository.New(o.Config, o.Log, store, agent)
	pagesRepository := pagesrepository.New(o.Log, templates, o.Config, store, agent)
	usersRepository := usersrepository.New(o.Log, o.Config, store)

	// Init Handlers
	expHandler := expressions.New(*orchRepository)
	pageHandler := pages.New(*pagesRepository)
	usersHandler := users.New(*usersRepository)

	// Router Initiziling
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	var once sync.Once

	// Route Handlers
	router.Route("/", func(r chi.Router) {
		r.Get("/", pageHandler.GetAuthPage)        // Get Auth Page
		r.Get("/register", pageHandler.GetRegPage) // Get Reg Page
		r.Post("/register", usersHandler.PostUser) // Post New User
		r.Get("/login", pageHandler.GetLoginPage)  // Get Login Page
		r.Post("/login", usersHandler.PostLogin)   // Post New User
		r.Get("/{user}", func(w http.ResponseWriter, r *http.Request) {
			user := chi.URLParam(r, "user")
			pageHandler.GetMainPage(w, r, user)
			// Check new expressions, parsing and calculate
			once.Do(func() {
				done := make(chan struct{})
				go orchRepository.Processing(user, o.Log, o.Config.ProcessingInterval, done)
			})
		}) // Get Main Page
		r.Post("/{user}", expHandler.PostExpression)             // Add Expression
		r.Get("/settings", pageHandler.GetSettingsPage)          // Get Settings Page
		r.Post("/settings", pageHandler.SetSettingsPage)         // Post Settings Page
		r.Get("/expressions/{user}", pageHandler.GetExpressions) // Get All Expressions
		r.Get("/tasks", pageHandler.GetTasks)                    // Get Tasks List
	})

	// Start HTTP Server
	o.Log.Info("starting server address", slog.String("address", o.Config.Address))

	o.Server = &http.Server{
		Addr:         o.Config.Address,
		Handler:      router,
		ReadTimeout:  o.Config.HTTPServer.Timeout,
		WriteTimeout: o.Config.HTTPServer.Timeout,
		IdleTimeout:  o.Config.HTTPServer.IdleTimeout,
	}

	if err := o.Server.ListenAndServe(); err != nil {
		o.Log.Info("failed to start server")
		return err
	}

	o.Log.Error("server stopped")

	return fmt.Errorf("server is stopped")
}

func (o *Orchestrator) Stop(ctx context.Context) error {

	err := o.Server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
