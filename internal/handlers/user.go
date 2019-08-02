package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gidyon/api-open-banking/internal/models/user"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func registerUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &user.User{}

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
	appUser := &user.User{}

	err := json.NewDecoder(r.Body).Decode(appUser)
	if err != nil {
		http.Error(w, "failed to decode request data", http.StatusInternalServerError)
		return
	}

	// Call Login method to authenticate user
	token, err := user.Login(appUser.UserID, appUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write token to response
	fmt.Fprint(w, token)
}

func getUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")

	appUser, err := user.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(appUser)
	if err != nil {
		http.Error(w, "failed to encode user data", http.StatusInternalServerError)
		return
	}
}

func updateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")

	newUser := &user.User{}

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, "failed to decode request data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = user.UpdateUser(userID, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
