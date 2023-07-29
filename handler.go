package main

import (
	"context"
	"fast_chat/core/port"
)

type FastChatServiceImp struct {
	httpHandler      port.HttpHandler
	websocketHandler port.WebSocketHandler
}

func (f *FastChatServiceImp) Run(ctx context.Context) error {
	if err := f.websocketHandler.Run(ctx); err != nil {
		panic(err)
	}
	if err := f.httpHandler.Run(ctx); err != nil {
		panic(err)
	}
	return nil
}
