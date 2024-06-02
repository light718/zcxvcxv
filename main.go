package main

import (
	"github/SEq7GVODOU4NEdrJu/websocketserver/pkg"
	"log"
	"net/http"
)

func main() {
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
