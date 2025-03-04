package router

import (
	"encoding/json"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"net/http"
)

func (r *Router) RegisterHandler(writer http.ResponseWriter, request *http.Request) {

	const nf = "register handler"

	var u models.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&u)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
		http.Error(writer, "bad request", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	fmt.Println(u)

	err = r.userService.Register(request.Context(), &u)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
	}
}
