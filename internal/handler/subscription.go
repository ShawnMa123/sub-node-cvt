package handler

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/ShawnMa123/sub-node-cvt/internal/converter"
)

// SubscriptionHandler 处理 /sub 路径的请求
func SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 从 URL query 获取参数
	query := r.URL.Query()
	nodesB64 := query.Get("nodes")
	rules := query.Get("rules") // e.g., "adguard,gfw"
	chainsB64 := query.Get("chains")

	if nodesB64 == "" {
		http.Error(w, "parameter 'nodes' is required", http.StatusBadRequest)
		return
	}

	// 2. Base64 解码
	nodesYAML, err := base64.URLEncoding.DecodeString(nodesB64)
	if err != nil {
		http.Error(w, "invalid base64 for 'nodes'", http.StatusBadRequest)
		return
	}

	var chainsJSON []byte
	if chainsB64 != "" {
		chainsJSON, err = base64.URLEncoding.DecodeString(chainsB64)
		if err != nil {
			http.Error(w, "invalid base64 for 'chains'", http.StatusBadRequest)
			return
		}
	}

	// 3. 调用核心逻辑生成配置
	finalConfig, err := converter.GenerateConfig(string(nodesYAML), rules, string(chainsJSON))
	if err != nil {
		log.Printf("Error generating config: %v", err)
		http.Error(w, "failed to generate config", http.StatusInternalServerError)
		return
	}

	// 4. 返回最终的配置文件
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(finalConfig)
}