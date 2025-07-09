package main

import (
	"log"
	"net/http"

	"github.com/ShawnMa123/sub-node-cvt/internal/handler"
)

func main() {
	// 注册 HTTP handler
	http.HandleFunc("/sub", handler.SubscriptionHandler)

	// 定义服务器端口
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	log.Printf("Access URL: http://localhost:%s/sub?...", port)

	// 启动服务器
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
