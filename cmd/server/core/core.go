package core

import (
	"github.com/TechCatsLab/apix/http/server"
)

func WriteStatusAndDataJSON(ctx *server.Context, status int, data interface{}) error {
	if data == nil {
		return ctx.ServeJSON(map[string]interface{}{"status": status})
	}

	return ctx.ServeJSON(map[string]interface{}{
		"status": status,
		"data":   data,
	})
}

func WriteStatusAndIDJSON(ctx *server.Context, status int, id interface{}) error {
	return ctx.ServeJSON(map[string]interface{}{
		"status": status,
		"ID":     id,
	})
}
