package router

import (
	"encoding/json"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/jwtc"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (r *Router) LoginHandler(writer http.ResponseWriter, request *http.Request) {
	const nf = "login handler"

	var u models.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&u)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
		http.Error(writer, "bad request", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	token, err := r.service.Login(request.Context(), &u)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))

		if errors.Is(err, store.ErrNotAuth) {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(writer, "error saved", http.StatusInternalServerError)
		return
	}
	logger.Log.Info(nf, fmt.Sprintf(" token: %v", token))

	cookie := &http.Cookie{
		Name:     jwtc.AuthCookie,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(writer, cookie)
	writer.WriteHeader(http.StatusOK)
}
