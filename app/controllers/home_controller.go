package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

func Home(rw http.ResponseWriter, r *http.Request) {

	render := render.New(render.Options{
		Layout: "layout",
	})

	_ = render.HTML(rw, http.StatusOK, "home", map[string]interface{}{
		"title": "Home Title",
		"body":  "Home Description",
	})
}
