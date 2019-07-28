package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/api-open-banking/internal/handler"
	"net/http"
)

func main() {
	port := ":8080"

	mux := handler.CreateMuxer()

	logrus.Printf("server started on port: %v\n", port)
	logrus.Fatalln(http.ListenAndServe(":8080", mux))
}
