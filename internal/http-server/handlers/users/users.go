package users

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/usersrepository"
)

type UsersHandle struct {
	UsersRepository usersrepository.UsersRepository
}

type Request struct {
}

type Response struct {
	response.Response
	Method string
}

type UserSaver interface {
	UserSave(string) (int, error)
}

func New(usersrep usersrepository.UsersRepository) *UsersHandle {
	return &UsersHandle{
		UsersRepository: usersrep,
	}
}

func (u *UsersHandle) PostUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.useres.saveuser"

	u.UsersRepository.Log = u.UsersRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
}
