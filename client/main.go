package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func test() {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	// 读取消息
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Received message:", string(p), messageType)
}

func test1() int {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err:%v\n", err)
			break
		}

		data = strings.TrimSpace(data)

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Fatal(err)
		}

		// 读取消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("received: ", string(p), messageType)
	}

	return 0
}

func main() {
	test1()
}
