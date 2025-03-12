package router

import (
	"encoding/json"
	"fmt"
	"github.com/EClaesson/go-luhn"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func (r *Router) OrderHandler(writer http.ResponseWriter, request *http.Request) {
	const nf = "order handler"

	/*	200 — номер заказа уже был загружен этим пользователем;
		202 — новый номер заказа принят в обработку;
		400 — неверный формат запроса;
		401 — пользователь не аутентифицирован;
		409 — номер заказа уже был загружен другим пользователем;
		422 — неверный формат номера заказа;
		500 — внутренняя ошибка сервера.
	*/

	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	ok, err := luhn.IsValid(string(data))
	if err != nil || !ok {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctx := request.Context()
	logger.Log.Debug(nf, "get userId from context: ", ctx.Value(UserId))

	id, ok := ctx.Value(UserId).(int64)
	if !ok {
		logger.Log.Error(nf, fmt.Sprintf("user id = %v not get in context", id))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.service.SaveOrder(ctx, string(data), id)
	if err != nil {
		if errors.Is(err, store.Busy) {
			writer.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, store.Double) {
			writer.WriteHeader(http.StatusOK)
			return
		}
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	return
}

func (r *Router) OrderList(writer http.ResponseWriter, request *http.Request) {
	const nf = "order list"

	/*
		204 — нет данных для ответа.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/

	result, err := r.service.OrderList(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(result) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	encoder := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = encoder.Encode(result); err != nil {
		http.Error(writer, "error write json", http.StatusInternalServerError)
		return
	}
	return
}
