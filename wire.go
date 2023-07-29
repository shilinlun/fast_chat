//go:generate wire
//go:build wireinject
// +build wireinject

package main

import (
	"fast_chat/core/port"
	"fast_chat/core/service/http_handler"
	"fast_chat/core/service/websocket_handler"
	"fast_chat/infra/storage"
	"fast_chat/infra/websocket_service"
	"github.com/google/wire"
)

func injectWebsocketHandler() port.WebSocketHandler {
	panic(
		wire.Build(
			wire.Bind(new(port.Storage), new(*storage.Mongo)),
			storage.ProvideMongo,
			wire.Bind(new(port.WebSocket), new(*websocket_service.WebsocketService)),
			websocket_service.ProvideWebsocketService,
			wire.Bind(new(port.WebSocketHandler), new(*websocket_handler.Handler)),
			websocket_handler.ProvideWebsocketHandler,
		),
	)
}

func injectHttpHandler() port.HttpHandler {
	panic(
		wire.Build(
			wire.Bind(new(port.HttpHandler), new(*http_handler.Handler)),
			http_handler.ProvideHandler,
		),
	)
}
