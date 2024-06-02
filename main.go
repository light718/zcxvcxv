package main

import (
	"github/SEq7GVODOU4NEdrJu/websocketserver/pkg"
	"log"
	"net/http"
	"runtime"
)

func main() {
	workerPool := pkg.NewWorkerPool(runtime.NumCPU())
	for i := 0; i < 100; i++ {
		workerPool.Submit(func(args ...interface{}) {

		}, i)
	}

	server := pkg.NewWebSocketServer()
	go server.Run()

	http.HandleFunc("/ws", server.HandleConnections)
	port := ":8080"
	log.Printf("服务器启动在端口 %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
