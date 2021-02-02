package main

import (
	"fmt"
	"net/http"
	"time"

	"dql_admin_backend/utils"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *utils.Connection
		token  []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = utils.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("{\"fff\":33}")); err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		if token, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		// 这里可根据token判断对方用户的身份
		fmt.Println("token:", string(token))
	}

ERR:
	fmt.Println("closing~~~")
	conn.Close()

}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("listening port: 7777")
	http.ListenAndServe("0.0.0.0:7777", nil)
}
