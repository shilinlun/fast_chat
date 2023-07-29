package websocket_handler

import (
	"context"
	"fast_chat/core/port"
)

type Handler struct {
	websocket port.WebSocket
}

func ProvideWebsocketHandler(
	websocket port.WebSocket,
) *Handler {
	return &Handler{
		websocket: websocket,
	}
}

func (h *Handler) Run(ctx context.Context) error {
	if err := h.websocket.Run(ctx); err != nil {
		panic(err)
	}
	return nil
}
