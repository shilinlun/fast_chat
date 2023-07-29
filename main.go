package main

import "context"

func main() {
	fastChat := &FastChatServiceImp{
		httpHandler:      injectHttpHandler(),
		websocketHandler: injectWebsocketHandler(),
	}
	if err := fastChat.Run(context.Background()); err != nil {
		panic(err)
	}
}
