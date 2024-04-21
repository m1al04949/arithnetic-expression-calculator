package users

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/lib/response"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/repositories/usersrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/pkg/hash"
)

type UsersHandle struct {
	UsersRepository usersrepository.UsersRepository
}

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New(usersrep usersrepository.UsersRepository) *UsersHandle {
	return &UsersHandle{
		UsersRepository: usersrep,
	}
}

// User Registration
func (u *UsersHandle) PostUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.useres.postuser"

	u.UsersRepository.Log = u.UsersRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

	req.Login = r.FormValue("login")
	req.Password = r.FormValue("password")
	// if API request
	if req.Login == "" && req.Password == "" {
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			u.UsersRepository.Log.Error("failed to decode request body", slog.Any("error", err.Error()))
			render.JSON(w, r, response.ErrorRequest("failed to decode request"))
			return
		}
	}

	u.UsersRepository.Log.Info("request decoded", slog.Any("request", req))

	//Check username on DB
	exists, err := u.UsersRepository.CheckUserOnDb(req.Login)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}
	if exists {
		render.JSON(w, r, response.ErrorRegistration("user is exists"))
		return
	}

	//Hashing password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("error hashing password"))
		return
	}

	// Create User
	err = u.UsersRepository.CreateUser(req.Login, hashedPassword)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("error creating user"))
		return
	}

	render.JSON(w, r, response.OK())
}

// User Authorization
func (u *UsersHandle) PostLogin(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.useres.postlogin"

	u.UsersRepository.Log = u.UsersRepository.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

	req.Login = r.FormValue("login")
	req.Password = r.FormValue("password")

	if req.Login == "" && req.Password == "" {
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			u.UsersRepository.Log.Error("failed to decode request body", slog.Any("error", err.Error()))
			render.JSON(w, r, response.ErrorRequest("failed to decode request"))
			return
		}
	}

	u.UsersRepository.Log.Info("request decoded", slog.Any("request", req))

	// Check User and Password
	check, err := u.UsersRepository.CheckAuthorization(req.Login, req.Password)
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}
	if !check {
		render.JSON(w, r, response.ErrorAuthorization("user is unauthorize"))
		return
	}

	// Create JWT токен
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = req.Login
	claims["exp"] = time.Now().Add(u.UsersRepository.Config.TokenExpire).Unix()

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(u.UsersRepository.Config.Token))
	if err != nil {
		render.JSON(w, r, response.ErrorServer("server internal error"))
		return
	}

	// Set cookie with token
	expiration := time.Now().Add(u.UsersRepository.Config.TokenExpire)
	cookie := http.Cookie{
		Name:     "jwtToken",
		Value:    tokenString,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	render.JSON(w, r, response.Authorization(req.Login, tokenString))
}
