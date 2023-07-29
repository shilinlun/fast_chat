package websocket_service

import (
	"context"
	"encoding/json"
	"fast_chat/core/entity"
	"fast_chat/core/port"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebsocketService struct {
	storage port.Storage
}

func ProvideWebsocketService(
	storage port.Storage,
) *WebsocketService {
	return &WebsocketService{
		storage: storage,
	}
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *WebsocketService) handleWebsocket(h http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(h, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	ctx := context.Background()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("received: %s", msg)
		if err := w.handleMsg(ctx, string(msg)); err != nil {
			err = conn.WriteMessage(messageType, []byte("failed"))
			if err != nil {
				log.Println(err)
				return
			}
			continue
		}
		if err := conn.WriteMessage(messageType, []byte("success")); err != nil {
			log.Println(err)
			return
		}
	}
}

func (w *WebsocketService) handleMsg(ctx context.Context, msg string) error {
	innerMsg := &entity.Msg{}
	err := json.Unmarshal([]byte(msg), innerMsg)
	if err != nil {
		return err
	}
	fmt.Println(innerMsg)
	if err := w.storage.Insert(ctx, innerMsg); err != nil {
		return err
	}
	return nil
}

func (w *WebsocketService) Run(ctx context.Context) error {
	http.HandleFunc("/ws", w.handleWebsocket)
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			panic(err)
		}
	}()
	return nil
}
