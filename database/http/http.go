package http

import (
	"context"
	"io"

	"trpc.group/trpc-go/trpc-go/client"
)

// ReqHeader 自定义请求字段
type ReqHeader struct {
	Header map[string]string // http的header头
	Path   string            // 默认空, 使用配置的Path; 如果填写更换path
	Query  string            // 直接拼接到path后面, 格式"?a=1&b=2"; 注意?以及encode
	Body   io.Reader         // 支持自定义格式
}

type HttpProxy interface {
	// Get 默认rsp内是json结构体, 如果非json, 可以从opts内传入client.WithReqHead(header)获取原始内容
	Get(ctx context.Context, rspbody interface{}, reqHeader *ReqHeader, opts ...client.Option) error

	// Post 默认body是json结构, 如果传入非json, 通过构造opts的client.WithReqHead替换
	Post(ctx context.Context, reqbody interface{}, rspbody interface{},
		reqHeader *ReqHeader, opts ...client.Option) error
	GetUrl() string
}

// Get 获取gorm, 不存在直接报错
func Get(name string) HttpProxy {
	return imp.get(name)
}
