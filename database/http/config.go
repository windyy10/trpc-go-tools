package http

// ("dsn://{user}:{password}@tcp({ip}:{port})/{db}?timeout=1s&parseTime=true"),
type Config struct {
	Providers []struct {
		Name    string `yaml:"name"`    // 名称
		Service string `yaml:"service"` // 业务名称, 自定义, 对应trpc的service流程, 格式: trpc.mysql.{module}.{db}
		Url     string `yaml:"url"`
		Target  string `yaml:"target"`  // 可选, 配置target强制使用target寻址; 无target直接用url解析的域名寻址
		Timeout int    `yaml:"timeout"` // 超时时间, ms
		Headers []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		} `yaml:"headers"`
		MaxAttempts int `yaml:"max_attempts"` // 可选, 重试次数
	} `yaml:"providers,omitempty"` // 列表
}
