//
//
//  wpim
//
//  Created by 甘文鹏 on 2018/5/27.
//  Copyright (c) 2018 甘文鹏. All rights reserved.
//

package socket

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
)

func newSocket() *Socket {
	ss := new(Socket)
	ss.connectPool = make(map[string]*websocket.Conn)
	return ss
}

type Socket struct {
	connectPool map[string]*websocket.Conn
}

/*接收请求*/
func (ss *Socket) ServeHTTP(respose http.ResponseWriter, req *http.Request) {
	// 身份校验
	token := req.URL.Query().Get("token")
	logs.Debug("TOKEN：%+v", token)

	if valid, msg := ss.checkToken(token); false==valid {
		logs.Debug(msg)
		respose.Write([]byte(msg))
		return
	}

	// HTTP请求升级为WebSocket
	conn, err := websocket.Upgrade(respose, req, nil, 1024, 1024)
	// 错误处理
	if err != nil {
		logs.Error(err)
		return
	}

	// 将连接保存到连接池
	ss.connectPool[token] = conn

	// 设置错误回调
	conn.SetCloseHandler(func(code int, text string) error {
		delete(ss.connectPool, token)
		return nil
	})

	// 监听连接
	ss.listenConn(conn)
}

/*监听Socket连接*/
func (ss *Socket) listenConn(conn *websocket.Conn) {
	for {
		// 从Socket通道读取
		typee, data, err := conn.ReadMessage()

		// 错误处理
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				break
			}

			logs.Error(err)
			break
		}

		// 判空
		if 0 == len(data) {
			continue
		}

		// 按照消息类型分发
		switch typee {
		case websocket.TextMessage:
			ss.readMessage(conn, []byte(data))
		case websocket.BinaryMessage:
			ss.readMessage(conn, data)
		}

		logs.Debug("%v", typee)
	}
}

/*读取消息*/
func (ss *Socket) readMessage(conn *websocket.Conn, msg []byte) {
	logs.Debug("%+v", string(msg))

	conn.WriteMessage(websocket.TextMessage, []byte("123"))
}

/*校验TOKEN合法性，true-合法，false-非法*/
func (ss *Socket) checkToken(token string) (valid bool, msg string){
	if 0 == len(token) {
		return false, "身份信息不得为空"
	}

	return true, ""
}
