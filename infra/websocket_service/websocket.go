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
	"sync"
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

var (
	// 升级成 WebSocket 协议
	upGrader = websocket.Upgrader{
		// 允许CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn        *websocket.Conn
	err         error
	handler     CenterHandler
	serviceOnce sync.Once
)

// CenterHandler 处理中心，关联着每个 Client 的注册、注销、广播通道，相当于每个用户的中心通讯的中介。
type CenterHandler struct {
	// 单人
	single chan []byte
	// 广播通道，有数据则循环每个用户广播出去
	broadcast chan []byte
	// 注销通道，有用户关闭连接 则将该用户剔出集合map中
	unregister chan *Client
	// 用户集合，每个用户本身也在跑两个协程，监听用户的读、写的状态
	clients map[string]*Client
}

// Client 抽象出来的 Client，里面有这个 websocket 连接的 读 和 写 操作
type Client struct {
	handler *CenterHandler
	conn    *websocket.Conn
	// 每个用户自己的循环跑起来的状态监控
	send    chan []byte
	storage port.Storage
}

// 写，主动推送消息给客户端
func (c *Client) writePump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 广播推过来的新消息，马上通过websocket推给自己
		message, _ := <-c.send
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 读，监听客户端是否有推送内容过来服务端
func (c *Client) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 循环监听是否该用户是否要发言
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 异常关闭的处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message []byte) {
	innerMSg := &entity.Msg{}
	if err := json.Unmarshal(message, innerMSg); err != nil {
		return
	}
	fmt.Println("start handleMessage", innerMSg)
	go func() {
		if err := c.storage.Insert(context.Background(), innerMSg); err != nil {
			return
		}
	}()
	switch innerMSg.MsgType {
	case 1:
		c.handler.clients[innerMSg.FromId] = c
	case 2:
		delete(c.handler.clients, innerMSg.FromId)
	case 3:
		c.send <- message
		if conn, ok := c.handler.clients[innerMSg.ToId]; ok {
			conn.send <- message
		}
	case 4:
		for _, client := range c.handler.clients {
			client.send <- message
		}
	}
	fmt.Println("end handleMessage", innerMSg)
}

func (w *WebsocketService) handleWebsocket(writer http.ResponseWriter, request *http.Request) {
	serviceOnce.Do(func() {
		// 应用一运行，就初始化 CenterHandler 处理中心对象
		handler = CenterHandler{
			single:     make(chan []byte),
			broadcast:  make(chan []byte),
			unregister: make(chan *Client),
			clients:    make(map[string]*Client),
		}

	})
	// 由 http 升级成为 websocket 服务
	if conn, err = upGrader.Upgrade(writer, request, nil); err != nil {
		log.Println(err)
		return
	}
	// 为每个连接创建一个 Client 实例，（实际上这里应该还有绑定用户真实信息的操作）
	client := &Client{&handler, conn, make(chan []byte, 256), w.storage}
	go client.writePump()
	go client.readPump()
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
