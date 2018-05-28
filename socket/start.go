//
//  
//  wpim
// 
//  Created by 甘文鹏 on 2018/5/27.
//  Copyright (c) 2018 甘文鹏. All rights reserved.
//  

package socket

import (
	"net/http"
	"github.com/astaxie/beego/logs"
)

func Start() {
	ss := newSocket()
	logs.Info("WebSocket is running ...")
	http.ListenAndServe(":8081", ss)
}
