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
	"time"
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
	// 用户集合，每个用户本身也在跑两个协程，监听用户的读、写的状态
	clients      map[string]*Client
	groupClients map[string][]string
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
		c.conn.Close()
	}()
	for {
		message, _ := <-c.send
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 读，监听客户端是否有推送内容过来服务端
func (c *Client) readPump() {
	defer func() {
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
	innerMSg.CreatTime = time.Now().Unix()
	fmt.Println("start handleMessage", innerMSg)
	go func() {
		if err := c.storage.Insert(context.Background(), innerMSg); err != nil {
			return
		}
	}()
	switch innerMSg.MsgType {
	case 1:
		c.handLoginMsg(innerMSg)
	case 2:
		c.handLogoutMsg(innerMSg)
	case 3:
		c.handSingleMsg(innerMSg, message)
	case 4:
		c.handAllMsg(message)
	case 5:
		c.handGroupMsg(innerMSg, message)
	}
	fmt.Println("end handleMessage", innerMSg)
}

// 登录
func (c *Client) handLoginMsg(innerMSg *entity.Msg) {
	c.handler.clients[innerMSg.FromId] = c
	var groupIds []string
	// 获取该uid所在的group ids
	for _, groupId := range groupIds {
		c.handler.groupClients[groupId] = append(c.handler.groupClients[groupId], innerMSg.FromId)
	}
}

// 注销
func (c *Client) handLogoutMsg(innerMSg *entity.Msg) {
	delete(c.handler.clients, innerMSg.FromId)
}

// 单对单
func (c *Client) handSingleMsg(innerMSg *entity.Msg, message []byte) {
	c.send <- message
	if conn, ok := c.handler.clients[innerMSg.ToId]; ok {
		conn.send <- message
	}
}

// 对所有人
func (c *Client) handAllMsg(message []byte) {
	for _, client := range c.handler.clients {
		client.send <- message
	}
}

// 群发
func (c *Client) handGroupMsg(innerMSg *entity.Msg, message []byte) {
	for _, toId := range c.handler.groupClients[innerMSg.ToId] {
		if conn, ok := c.handler.clients[toId]; ok {
			conn.send <- message
		}
	}
}

func (w *WebsocketService) handleWebsocket(writer http.ResponseWriter, request *http.Request) {
	serviceOnce.Do(func() {
		// 应用一运行，就初始化 CenterHandler 处理中心对象
		handler = CenterHandler{
			clients:      make(map[string]*Client),
			groupClients: make(map[string][]string),
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
