package handler

import (
	"demo/constants"
	"demo/core"
	"demo/models"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/logging/logrus"
)

// Show repos
func Show(c *server.Context) error {
	var (
		pkg  string
		data interface{}
	)
	pkg = c.Request().FormValue("pkg")
	data, err := models.Find(models.Deal(pkg))
	if err != nil {
		logrus.Error("Find():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, data)
}
