package pages

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/pagesrepository"
)

type PagesHandle struct {
	PagesRepository pagesrepository.PagesRepository
}

type Response struct {
	response.Response
	Method string
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

	data := "Main Page"

	err := h.PagesRepository.Templates.Main.Execute(w, data)
	if err != nil {
		h.PagesRepository.Log.Error("failed to download main page")
		render.JSON(w, r, response.ErrorRequest("failed to download main page"))
		return
	}
}
