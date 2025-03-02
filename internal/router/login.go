package router

import "net/http"

func (r *Router) LoginHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("LoginHandler"))
}
