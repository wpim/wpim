//
//  
//  wpim
// 
//  Created by 甘文鹏 on 2018/5/27.
//  Copyright (c) 2018 甘文鹏. All rights reserved.
//  

package socket

import "net/http"

func Start() {
	ss := newSocket()
	http.ListenAndServe(":8081", ss)
}
