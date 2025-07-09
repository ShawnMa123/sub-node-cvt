package handler

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/ShawnMa123/sub-node-cvt/internal/converter"
)

// SubscriptionHandler 只处理通过 URL 参数生成配置的请求
func SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nodesB64 := query.Get("nodes")

	if nodesB64 == "" {
		http.Error(w, "parameter 'nodes' is required", http.StatusBadRequest)
		return
	}

	nodesYAML, err := base64.RawURLEncoding.DecodeString(nodesB64)
	if err != nil {
		http.Error(w, "invalid base64 for 'nodes'", http.StatusBadRequest)
		return
	}

	rules := query.Get("rules")
	chainsB64 := query.Get("chains")

	var chainsJSON []byte
	if chainsB64 != "" {
		chainsJSON, err = base64.RawURLEncoding.DecodeString(chainsB64)
		if err != nil {
			http.Error(w, "invalid base64 for 'chains'", http.StatusBadRequest)
			return
		}
	}

	finalConfig, err := converter.GenerateConfig(string(nodesYAML), rules, string(chainsJSON))
	if err != nil {
		log.Printf("Error generating config: %v", err)
		http.Error(w, "failed to generate config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(finalConfig)
}
