package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// BASEURL is base URL for accessing the open banking api
const BASEURL = "https://apisandbox.openbankproject.com"

// CreateMuxer creates an http muxer
func CreateMuxer() http.Handler {
	// Register endpoints
	mux := httprouter.New()

	// CRUD operations on our API
	// Register user to our API
	mux.POST("/users/register", registerUserHandler)
	// Logins a user to our API
	mux.POST("/users/login", loginUserHandler)
	// Updates user profile in our API DB
	mux.PUT("/users/:userID", updateUserHandler)
	// Retrieves a user profile from our API DB
	mux.GET("/users/:userID", getUserHandler)

	// Working with Banks
	// Create a bank to the open-banking external API
	mux.POST("/users/banks/:bankID/create-customer", createBankCustomerHandlerV2)
	// Adds a bank to a list of user banks to our API
	mux.PUT("/users/:userID/banks/add/:bankID", addBankHandler)
	// Retrives list banks a user is customer
	mux.GET("/users/:userID/banks", getBanksHandler)
	// Retrieve a single bank a user is customer
	mux.GET("/users/:userID/banks/:bankID", getBankHandler)

	return mux
}
