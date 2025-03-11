package router

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/jwtc"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"net/http"
)

type keyUserID string

const UserId keyUserID = "userID"

func (r *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		logger.Log.Info("AuthMiddleware")
		cookie, err := request.Cookie(jwtc.AuthCookie)
		if err != nil {
			http.Error(writer, store.ErrNotAuth.Error(), http.StatusUnauthorized)
			return
		}

		id, err := jwtc.Decode(cookie.Value, r.Config)
		if err != nil {
			logger.Log.Error(err)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), UserId, id)
		request.WithContext(ctx)
		next.ServeHTTP(writer, request)
	})
}
