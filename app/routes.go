package app

import (
	"go_shop/app/controllers"
)

func (server *Server) InitializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
