package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/api-open-banking/internal/handler"
	"net/http"
)

func main() {
	port := flag.String("port", ":9090", "port for open banking api server")
	flag.Parse()

	mux := handler.CreateMuxer()

	logrus.Printf("open banking server started on port: %s\n", *port)
	logrus.Fatal(http.ListenAndServe(*port, mux))
}
