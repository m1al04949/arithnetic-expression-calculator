package pages

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/model"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/pagesrepository"
)

type PagesHandle struct {
	PagesRepository pagesrepository.PagesRepository
}

type Response struct {
	response.Response
	Method string
}

type RequestSettings struct {
	OperationSum time.Duration
	OperationSub time.Duration
	OperationMul time.Duration
	OperationDiv time.Duration
}

type RequestExpressions struct {
	Expressions []model.ExpressionTab
}

func New(pagesrep pagesrepository.PagesRepository) *PagesHandle {
	return &PagesHandle{
		PagesRepository: pagesrep,
	}
}

func (h *PagesHandle) GetMainPage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getmainpage"

	h.PagesRepository.Log = h.PagesRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	data := "Главная страница"

	err := h.PagesRepository.Templates.Main.Execute(w, data)
	if err != nil {
		h.PagesRepository.Log.Error("failed to download main page")
		render.JSON(w, r, response.ErrorRequest("failed to download main page"))
		return
	}
}

func (h *PagesHandle) GetSettingsPage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getsettingspage"

	h.PagesRepository.Log = h.PagesRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	// data := "Settings Page"

	err := h.PagesRepository.Templates.Settings.Execute(w, h.PagesRepository.Config.Timeouts)
	if err != nil {
		h.PagesRepository.Log.Error("failed to download settings page")
		render.JSON(w, r, response.ErrorRequest("failed to download settings page"))
		return
	}
}

func (h *PagesHandle) SetSettingsPage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.setsettingspage"

	h.PagesRepository.Log = h.PagesRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var (
		req RequestSettings
		err error
	)

	sum := r.FormValue("addition_time")
	req.OperationSum, err = time.ParseDuration(sum)
	if err != nil {
		h.PagesRepository.Log.Error("failed to update settings operation sum")
		render.JSON(w, r, response.ErrorRequest("failed to update operation sum:"+err.Error()))
		return
	}
	sub := r.FormValue("subtraction_time")
	req.OperationSub, err = time.ParseDuration(sub)
	if err != nil {
		h.PagesRepository.Log.Error("failed to update settings operation sub")
		render.JSON(w, r, response.ErrorRequest("failed to update operation sub:"+err.Error()))
		return
	}
	mul := r.FormValue("multiplication_time")
	req.OperationMul, err = time.ParseDuration(mul)
	if err != nil {
		h.PagesRepository.Log.Error("failed to update settings operation mul")
		render.JSON(w, r, response.ErrorRequest("failed to update operation mul:"+err.Error()))
		return
	}
	div := r.FormValue("division_time")
	req.OperationDiv, err = time.ParseDuration(div)
	if err != nil {
		h.PagesRepository.Log.Error("failed to update settings operation div")
		render.JSON(w, r, response.ErrorRequest("failed to update operation div:"+err.Error()))
		return
	}

	h.PagesRepository.Config.Timeouts.OperationSumInterval = req.OperationSum
	h.PagesRepository.Config.Timeouts.OperationSubInterval = req.OperationSub
	h.PagesRepository.Config.Timeouts.OperationMulInterval = req.OperationMul
	h.PagesRepository.Config.Timeouts.OperationDivInterval = req.OperationDiv

	// data := "Settings Page"

	err = h.PagesRepository.Templates.Settings.Execute(w, h.PagesRepository.Config.Timeouts)
	if err != nil {
		h.PagesRepository.Log.Error("failed to download settings page")
		render.JSON(w, r, response.ErrorRequest("failed to download settings page"))
		return
	}
}

func (h *PagesHandle) GetExpressions(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.expressionpage"

	h.PagesRepository.Log = h.PagesRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	expressions, err := h.PagesRepository.Store.GetAllExpressions()
	if err != nil {
		render.JSON(w, r, response.ErrorRequest("failed to download expressions page"+err.Error()))
	}

	data := RequestExpressions{
		Expressions: expressions,
	}

	err = h.PagesRepository.Templates.Expressions.Execute(w, data)
	if err != nil {
		h.PagesRepository.Log.Error("failed to download expressions page")
		render.JSON(w, r, response.ErrorRequest("failed to download expressions page"+err.Error()))
		return
	}
}
