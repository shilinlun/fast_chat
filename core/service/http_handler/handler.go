package http_handler

import (
	"context"
	"fast_chat/router_handler/ping"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func ProvideHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Run(ctx context.Context) error {
	r := gin.Default()
	if err := startGin(r); err != nil {
		panic(err)
	}
	return nil
}

func startGin(r *gin.Engine) error {
	r.GET("/ping", ping.Ping)
	return r.Run()
}
