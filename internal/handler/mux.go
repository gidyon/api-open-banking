package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// CreateMuxer creates an http muxer
func CreateMuxer() http.Handler {
	// Register endpoints
	mux := httprouter.New()

	mux.POST("/user/register", registerUserHandler)
	mux.POST("/user/login", loginUserHandler)

	return mux
}
