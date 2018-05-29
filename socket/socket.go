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
	"sync"
)

func newSocket() *Socket {
	ss := new(Socket)
	ss.connectPool = make(map[string]*websocket.Conn)
	ss.lock = new(sync.Mutex)
	return ss
}

const paramToken = "im_token"

type Socket struct {
	connectPool map[string]*websocket.Conn
	lock *sync.Mutex
}

/*接收请求*/
func (ss *Socket) ServeHTTP(respose http.ResponseWriter, req *http.Request) {
	// 请求合法性校验
	if valid, msg := ss.validRequest(req); false==valid {
		logs.Debug(msg)
		respose.Write([]byte(msg))
		return
	}

	token := req.URL.Query().Get(paramToken)
	logs.Debug("TOKEN：%+v", token)

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

	//go ss.broadcast(msg)
}

/*广播*/
func (ss *Socket)broadcast(input []byte)  {
	ss.lock.Lock()

	// todo 广播
	for _, v := range ss.connectPool {
		v.WriteMessage(websocket.TextMessage, input)
	}

	ss.lock.Unlock()
}

/*校验请求合法性*/
func (ss *Socket) validRequest(req *http.Request) (valid bool, msg string){
	logs.Debug("%+v", req)

	token := req.URL.Query().Get(paramToken)
	if 0 == len(token) {
		return false, "身份信息不得为空"
	}

	return true, ""
}
