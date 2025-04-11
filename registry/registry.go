package registry

import (
	"sync"

	"gopkg.in/yaml.v3"
)

// Decoder 通用解析器, yaml没有对应interface, 解析器没法兼容
type Decoder interface {
	Decode(cfg any) error
}

var (
	setupFuncMap = make(map[string]func(decoder Decoder) error)
	lock         sync.Mutex
)

// Setup 初始化所有插件
func Setup(decoder Decoder) error {
	data := make(map[string]yaml.Node)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	for k, node := range data {
		if f, ok := setupFuncMap[k]; ok {
			if err := f(&node); err != nil {
				return err
			}
		}
	}
	return nil
}

// Registry decoder为yaml的配置解析, 提供给插件注册
func Registry(name string, setup func(decoder Decoder) error) {
	lock.Lock()
	defer lock.Unlock()
	setupFuncMap[name] = setup
}
