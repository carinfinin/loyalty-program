package router

import (
	"fmt"
	"github.com/EClaesson/go-luhn"
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
	
	rr, _ := r.service.SaveOrder(request.Context(), string(data))
	fmt.Println(rr)
}
