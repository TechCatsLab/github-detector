package main

import (
	"demo/conf"
	"demo/router"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/apix/http/server/middleware"
	"github.com/TechCatsLab/logging/logrus"
)

var (
	ep *server.Entrypoint
)

func start() {
	ep = server.NewEntrypoint(&server.Configuration{
		Address: conf.Config.Port,
	}, nil)

	ep.AttachMiddleware(middleware.NegroniRecoverHandler())
	ep.AttachMiddleware(middleware.NegroniLoggerHandler())
	ep.AttachMiddleware(middleware.NegroniCorsAllowAll())

	if err := ep.Start(router.Router.Handler()); err != nil {
		logrus.Error(err)
		return
	}

	ep.Wait()
}
