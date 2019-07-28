package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gidyon/api-open-banking/internal/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func registerUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "failed to decode request data", http.StatusInternalServerError)
		return
	}

	// Call Register method to register user
	userID, err := user.Register()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write some string to
	fmt.Fprintf(w, "registred successfully: id: %s", userID)
}

func loginUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "failed to decode request data", http.StatusInternalServerError)
		return
	}

	// Call Login method to authenticate user
	token, err := user.Login()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write token to response
	fmt.Fprint(w, token)
}
