package main

import (
	"fmt"
	"rtsp2hls/src/help"
	"rtsp2hls/src/routers"
)

func main() {
	addr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	router := routers.InitRouter()
	router.Run(addr)
}
