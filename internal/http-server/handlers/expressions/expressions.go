package expressions

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/orchrepository"
)

type ExpressionHandle struct {
	OrchRepository orchrepository.OrchRepository
}

type Request struct {
	User       string `json:"user"`
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

	cookie, err := r.Cookie("jwtToken")
	if err != nil {
		h.OrchRepository.Log.Error("no token in cookie", slog.Any("error", err.Error()))
		return
	}

	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.OrchRepository.Cfg.Token), nil
	})
	if err != nil {
		h.OrchRepository.Log.Error("failed parsing token", slog.Any("error", err.Error()))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["user"].(string)
		h.OrchRepository.Log.Info("user authentication", slog.Any("username", username))
	} else {
		h.OrchRepository.Log.Error("failed token")
		return
	}

	req.Expression = r.FormValue("expression")
	req.User = r.FormValue("user")

	h.OrchRepository.Log.Info("request body decoded", slog.Any("request", req))

	//Validate expression
	err, check := h.OrchRepository.CheckExpression(req.Expression)
	if !check {
		render.JSON(w, r, response.ErrorExpression("bad expression: "+err.Error()))
		return
	}

	//Check expression on DB
	check, err = h.OrchRepository.CheckExpOnDb(req.User, req.Expression)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}
	if !check {
		render.JSON(w, r, response.ErrorExpression("expression is exists"))
		return
	}

	//Add expression to DB
	id, err := h.OrchRepository.AddExpression(req.User, req.Expression)
	if err != nil {
		render.JSON(w, r, response.ErrorServer(err.Error()))
		return
	}

	h.OrchRepository.Log.Info("expression added", slog.Int("id", id))

	// Перенаправление на предыдущую страницу
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
