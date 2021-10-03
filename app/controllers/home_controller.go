package controllers

import (
	"fmt"
	"net/http"
)

func Home(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Welcome to the Go Shop Home Page")
}
