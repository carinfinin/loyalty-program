package router

import "net/http"

func (r *Router) RegisterHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("RegisterHandler"))
}
