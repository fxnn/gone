package main

import (
	"net/http"
	"github.com/fxnn/gone/handler"
)

func main() {
	http.ListenAndServe(":8080", handler.NewHandler())
}
