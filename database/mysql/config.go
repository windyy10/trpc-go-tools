package mysql

// Config 配置
type Config struct {
	Providers []struct {
		Name    string `yaml:"name"`    // 名称
		Service string `yaml:"service"` // 业务名称, 自定义, 对应trpc的service流程, 格式: trpc.mysql.{module}.{db}
		Timeout int    `yaml:"timeout"` // 超时时间, ms
		Target  string `yaml:"target"`  // db地址, 格式 dsn://{user}:{password}@tcp(ip:port)/{db}?timeout=1s&parseTime=true
		Table   string `yaml:"table"`   // 表名
	} `yaml:"providers,omitempty"` // 列表
}
