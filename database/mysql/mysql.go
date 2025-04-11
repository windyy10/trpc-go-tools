package mysql

import (
	"context"

	"gorm.io/gorm"
)

// Get 获取gorm, 不存在时返回nil, 线程安全的使用
// name: 注册的实体名称
// ctx: 上下文信息, 确保该方法返回的gorm对象是clone出来的, 每次Get出都是一个新的session
func Get(name string, ctx context.Context) *gorm.DB {
	return imp.get(name).WithContext(ctx)
}
