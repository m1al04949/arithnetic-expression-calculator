package expressions

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
)

type Request struct {
	ExpressionID int `json:"exp_id"`
}

type Response struct {
	response.Response
	ExpressionID int `json:"exp_id"`
	Method       string
}

type ExpSaver interface {
	SaveExpression(string) (int, error)
}

func PostExpression(log *slog.Logger, expSaver ExpSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.postexpression"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request")
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		// log.Info("expression added", slog.Int("id", id))

		// render.JSON(w, r, Response{
		// 	Response:     response.OK(),
		// 	ExpressionID: id,
		// 	Method:       r.Method,
		// })
	}
}
