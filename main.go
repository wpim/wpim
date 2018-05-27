//
//  
//  wpim
// 
//  Created by 甘文鹏 on 2018/5/27.
//  Copyright (c) 2018 甘文鹏. All rights reserved.
//  

package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/wpim/wpim/socket"
)

func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func main() {
	socket.Start()
}