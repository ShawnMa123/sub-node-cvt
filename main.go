package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ShawnMa123/sub-node-cvt/internal/handler"
)

func main() {
	// 检查必要的环境变量
	if os.Getenv("GITHUB_CLIENT_ID") == "" || os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		log.Fatal("Error: GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET environment variables must be set.")
	}
	handler.InitOAuth()

	// 注册 API 和认证路由
	http.HandleFunc("/sub", handler.SubscriptionHandler)
	http.HandleFunc("/auth/github", handler.HandleGitHubLogin)
	http.HandleFunc("/auth/github/callback", handler.HandleGitHubCallback)
	http.HandleFunc("/api/gist", handler.HandleCreateGist)
	http.HandleFunc("/api/user", handler.HandleUserInfo) // 新增用于获取用户信息的端点

	// 注册前端静态文件服务
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
