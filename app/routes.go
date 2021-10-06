package app

import (
	"go_shop/app/controllers"

	"github.com/gorilla/mux"
)

func (server *Server) InitializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
