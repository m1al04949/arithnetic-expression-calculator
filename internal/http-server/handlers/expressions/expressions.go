package expressions

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/orchrepository"
)

type ExpressionHandle struct {
	OrchRepository orchrepository.OrchRepository
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	response.Response
	ExpressionID int `json:"exp_id"`
	Method       string
}

type ExpSaver interface {
	ExpressionSave(string) (int, error)
}

func New(orchrep orchrepository.OrchRepository) *ExpressionHandle {
	return &ExpressionHandle{
		OrchRepository: orchrep,
	}
}

func (h *ExpressionHandle) PostExpression(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.postexpression"

	h.OrchRepository.Log = h.OrchRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		h.OrchRepository.Log.Error("failed to decode request body")
		render.JSON(w, r, response.ErrorRequest("failed to decode request"))
		return
	}

	h.OrchRepository.Log.Info("request body decoded", slog.Any("request", req))

	//Validate expression
	check, err := h.OrchRepository.CheckExpression(req.Expression)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}
	if !check {
		render.JSON(w, r, response.ErrorExpression("bad expression"))
		return
	}

	//Check expression on DB
	check, err = h.OrchRepository.CheckExpOnDb(req.Expression)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}
	if !check {
		render.JSON(w, r, response.ErrorExpression("expression is exists"))
		return
	}

	//Add expression to DB
	id, err := h.OrchRepository.AddExpression(req.Expression)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}

	h.OrchRepository.Log.Info("expression added", slog.Int("id", id))

	render.JSON(w, r, Response{
		Response:     response.OK(),
		ExpressionID: id,
		Method:       r.Method,
	})
}

// func PostExpression(log *slog.Logger, expSaver ExpSaver) http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handlers.postexpression"

// 		log = log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)

// 		var req Request

// 		err := render.DecodeJSON(r.Body, &req)
// 		if err != nil {
// 			log.Error("failed to decode request body")
// 			render.JSON(w, r, response.Error("failed to decode request"))
// 			return
// 		}

// 		log.Info("request body decoded", slog.Any("request", req))

// 		//Validate expression

// 		//Add expression to DB

// 		// log.Info("expression added", slog.Int("id", id))

// 		render.JSON(w, r, Response{
// 			Response: response.OK(),
// 			// ExpressionID: id,
// 			Method: r.Method,
// 		})
// 	}
// }
