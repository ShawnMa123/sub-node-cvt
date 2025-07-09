package types

// Chain 定义了中转链的结构，对应前端传入的JSON
type Chain struct {
	Relay   string `json:"relay"`   // 中转节点的名称
	Landing string `json:"landing"` // 落地节点的名称
}

// RuleSetContent 定义了从 rulesets/ 目录中读取的单个规则文件的结构
type RuleSetContent struct {
	RuleProviders map[string]RuleProvider `yaml:"rule-providers"`
	Rules         []string                `yaml:"rules"`
}

// ClashConfig 是最终生成的 Clash 配置文件的完整结构
type ClashConfig struct {
	Port          int                     `yaml:"port"`
	SocksPort     int                     `yaml:"socks-port"`
	RedirPort     int                     `yaml:"redir-port"`
	TProxyPort    int                     `yaml:"tproxy-port"`
	MixedPort     int                     `yaml:"mixed-port"`
	AllowLan      bool                    `yaml:"allow-lan"`
	Mode          string                  `yaml:"mode"`
	LogLevel      string                  `yaml:"log-level"`
	ExternalUI    string                  `yaml:"external-ui"`
	DNS           DNS                     `yaml:"dns"`
	Proxies       []map[string]any        `yaml:"proxies"`
	ProxyGroups   []ProxyGroup            `yaml:"proxy-groups"`
	RuleProviders map[string]RuleProvider `yaml:"rule-providers,omitempty"`
	Rules         []string                `yaml:"rules"`
	TUN           TUN                     `yaml:"tun,omitempty"`
}

type DNS struct {
	Enable            bool     `yaml:"enable"`
	IPv6              bool     `yaml:"ipv6"`
	Listen            string   `yaml:"listen"`
	EnhancedMode      string   `yaml:"enhanced-mode"`
	FakeIPRange       string   `yaml:"fake-ip-range"`
	DefaultNameserver []string `yaml:"default-nameserver"`
	Nameserver        []string `yaml:"nameserver"`
	Fallback          []string `yaml:"fallback"`
	FallbackFilter    struct {
		GeoIP  bool     `yaml:"geoip"`
		IPCIDR []string `yaml:"ipcidr"`
	} `yaml:"fallback-filter"`
}

type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

type RuleProvider struct {
	Type     string `yaml:"type"`
	Behavior string `yaml:"behavior"`
	URL      string `yaml:"url"`
	Path     string `yaml:"path"`
	Interval int    `yaml:"interval"`
}

type TUN struct {
	Enable        bool     `yaml:"enable"`
	Stack         string   `yaml:"stack"`
	DNSHijack     []string `yaml:"dns-hijack"`
	AutoRoute     bool     `yaml:"auto-route"`
	AutoDetectInt bool     `yaml:"auto-detect-interface"`
}
