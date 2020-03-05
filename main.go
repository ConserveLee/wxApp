package main

import (
	"app/routers"
	"github.com/imroc/req"
	"net/http"
	"time"
)

func main()  {
	SetConnPool()
	routers.InitServer()
}

func SetConnPool() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		MaxIdleConnsPerHost: 2000,
		// 无需设置MaxIdleConns
		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		// MaxIdleConns 默认是0，0表示不限制
	}
	req.SetClient(client)
	req.SetTimeout(5 * time.Second)
}