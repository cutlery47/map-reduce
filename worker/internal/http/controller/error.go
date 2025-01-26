package controller

import "net/http"

func handleErr(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("internal server error: " + err.Error()))
}
