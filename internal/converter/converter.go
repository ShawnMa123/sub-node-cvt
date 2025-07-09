package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// --- 这里是修正的地方 ---
	"github.com/ShawnMa123/sub-node-cvt/internal/types"
	"gopkg.in/yaml.v3"
)

const (
	templatePath = "./templates/meta_template.yaml"
	rulesetsDir  = "./rulesets/"
)

// GenerateConfig 是生成配置文件的核心函数
func GenerateConfig(nodesYAML, rulesStr, chainsJSON string) ([]byte, error) {
	// 1. 加载基础模板
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("reading template file: %w", err)
	}

	var config types.ClashConfig
	if err := yaml.Unmarshal(templateData, &config); err != nil {
		return nil, fmt.Errorf("unmarshaling template: %w", err)
	}

	// 2. 解析用户提供的节点
	var proxies []map[string]any
	if err := yaml.Unmarshal([]byte(nodesYAML), &proxies); err != nil {
		return nil, fmt.Errorf("unmarshaling user nodes: %w", err)
	}
	config.Proxies = proxies

	// 提取所有节点名称
	var allNodeNames []string
	for _, p := range proxies {
		if name, ok := p["name"].(string); ok {
			allNodeNames = append(allNodeNames, name)
		}
	}

	// 3. 解析中转链并创建 Relay 代理组
	var chains []types.Chain
	if chainsJSON != "" {
		if err := json.Unmarshal([]byte(chainsJSON), &chains); err != nil {
			return nil, fmt.Errorf("unmarshaling chains json: %w", err)
		}
	}

	var chainGroupNames []string
	for _, chain := range chains {
		groupName := fmt.Sprintf("CHAIN | %s -> %s", chain.Relay, chain.Landing)
		relayGroup := types.ProxyGroup{
			Name:    groupName,
			Type:    "relay",
			Proxies: []string{chain.Relay, chain.Landing},
		}
		config.ProxyGroups = append(config.ProxyGroups, relayGroup)
		chainGroupNames = append(chainGroupNames, groupName)
	}

	// 4. 创建 "PROXY" 选择组
	proxySelectGroup := types.ProxyGroup{
		Name:    "PROXY",
		Type:    "select",
		Proxies: append(chainGroupNames, allNodeNames...), // 把中转链放在前面
	}
	// 将 "PROXY" 组插入到所有组的最前面
	config.ProxyGroups = append([]types.ProxyGroup{proxySelectGroup}, config.ProxyGroups...)

	// 5. 加载并合并规则
	config.RuleProviders = make(map[string]types.RuleProvider)

	// 首先添加兜底规则
	finalRules := []string{
		"GEOIP,CN,DIRECT",
		"MATCH,PROXY",
	}

	selectedRules := strings.Split(rulesStr, ",")
	for _, ruleName := range selectedRules {
		ruleName = strings.TrimSpace(ruleName)
		if ruleName == "" {
			continue
		}

		ruleFilePath := filepath.Join(rulesetsDir, ruleName+".yaml")
		ruleData, err := os.ReadFile(ruleFilePath)
		if err != nil {
			continue
		}

		var ruleSet types.RuleSetContent
		if err := yaml.Unmarshal(ruleData, &ruleSet); err != nil {
			continue
		}

		for name, provider := range ruleSet.RuleProviders {
			config.RuleProviders[name] = provider
		}
		finalRules = append(ruleSet.Rules, finalRules...)
	}

	config.Rules = finalRules

	// 6. 序列化最终配置
	finalYAML, err := yaml.Marshal(&config)
	if err != nil {
		return nil, fmt.Errorf("marshaling final config: %w", err)
	}

	return finalYAML, nil
}
