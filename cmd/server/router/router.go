package router

import (
	"demo/handler"

	"github.com/TechCatsLab/apix/http/server"
)

var (
	Router *server.Router
)

func init() {
	Router = server.NewRouter()
	register(Router)
}

func register(r *server.Router) {
	r.Get("/detector/v1/show", handler.Show)
}
