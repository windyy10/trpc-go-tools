package http

import (
	"context"
	"net/url"
	"time"

	"github.com/windyy10/trpc-go-tools/database/registry"
	"trpc.group/trpc-go/trpc-go/client"
	thttp "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

var (
	imp *httpImp
)

type httpImp struct {
	httpMap map[string]HttpProxy
}

type httpProxyImp struct {
	u           *url.URL
	proxy       thttp.Client
	url         string
	path        string
	headers     map[string]string
	maxAttempts int
}

func (h *httpProxyImp) GetUrl() string {
	return h.url
}

func (h *httpProxyImp) Post(ctx context.Context, reqbody interface{}, rspbody interface{},
	reqHeader *ReqHeader, opts ...client.Option) error {

	header := &thttp.ClientReqHeader{
		Schema: h.u.Scheme,
		Host:   h.u.Host,
	}
	// 支持自定义的urlencoded编码
	if reqHeader != nil {
		header.ReqBody = reqHeader.Body
	}

	// 添加header, 按照默认值->规范值->自定义顺序添加
	for k, v := range h.headers {
		header.AddHeader(k, v)
	}

	if reqHeader != nil {
		for k, v := range reqHeader.Header {
			header.AddHeader(k, v)
		}
	}
	// 大小写不敏感, 默认添加json格式
	if header.Header.Get("Content-Type") == "" {
		header.AddHeader("Content-Type", "application/json; charset=utf-8")
	}
	path := h.buildPath(reqHeader)
	o := make([]client.Option, 0, 1+len(opts))
	o = append(o, client.WithReqHead(header))
	o = append(o, opts...)

	attempts := 0
	for {
		err := h.proxy.Post(ctx, path, reqbody, rspbody, o...)
		if err == nil {
			return nil
		}
		// 重试的次数不能超过最大尝试次数
		attempts++
		if attempts > h.maxAttempts {
			return err
		}
	}
}

func (h *httpProxyImp) Get(ctx context.Context,
	rspbody interface{},
	reqHeader *ReqHeader,
	opts ...client.Option,
) error {
	header := &thttp.ClientReqHeader{
		Schema: h.u.Scheme,
		Host:   h.u.Host,
	}
	// 添加header, 按照默认值->规范值->自定义顺序添加
	for k, v := range h.headers {
		header.AddHeader(k, v)
	}
	if reqHeader != nil {
		for k, v := range reqHeader.Header {
			header.AddHeader(k, v)
		}
	}
	path := h.buildPath(reqHeader)
	o := make([]client.Option, 0, 1+len(opts))
	o = append(o, client.WithReqHead(header))
	o = append(o, opts...)

	attempts := 0
	for {
		err := h.proxy.Get(ctx, path, rspbody, o...)
		if err == nil {
			return nil
		}
		// 重试的次数不能超过最大尝试次数
		attempts++
		if attempts > h.maxAttempts {
			return err
		}
	}
}

func (h *httpProxyImp) buildPath(reqHeader *ReqHeader) string {
	path := h.path
	if reqHeader != nil {
		if reqHeader.Path != "" {
			path = reqHeader.Path
		}
		if reqHeader.Query != "" {
			path += reqHeader.Query
		}
	}
	return path
}

func (h *httpImp) get(name string) HttpProxy {
	if p, ok := h.httpMap[name]; ok {
		return p
	}
	return nil
}

func (h *httpImp) setup(decoder registry.Decoder) error {
	conf := &Config{}
	decoder.Decode(conf)
	log.Debugf("setup conf %+v", conf)
	for _, item := range conf.Providers {
		p := &httpProxyImp{
			url:         item.Url,
			headers:     make(map[string]string),
			maxAttempts: item.MaxAttempts,
		}
		if err := p.parseUrl(); err != nil {
			log.Errorf("new gorm clinet fail %s", err.Error())
			return err
		}
		for _, kv := range item.Headers {
			p.headers[kv.Key] = kv.Value
		}
		opts := []client.Option{
			client.WithTimeout(time.Duration(item.Timeout) * time.Millisecond),
			client.WithCalleeMethod("http"),
		}

		if item.Target != "" {
			opts = append(opts, client.WithTarget(item.Target))
		} else {
			opts = append(opts, client.WithTarget("dns://"+p.u.Host))
		}
		p.proxy = thttp.NewClientProxy(item.Service, opts...)
		h.httpMap[item.Name] = p
		log.Infof("http register %+v ok", item)
	}
	return nil
}

func (h *httpProxyImp) parseUrl() error {
	var err error
	// 解析url字段
	if h.u, err = url.Parse(h.url); err != nil {
		return err
	}
	// path需要重新拼接, 解析后字段都拆碎了
	h.path = h.u.Path
	if h.u.RawQuery != "" {
		h.path += "?" + h.u.RawQuery
	}
	return nil
}

func init() {
	imp = &httpImp{
		httpMap: make(map[string]HttpProxy),
	}
	registry.Registry("http", imp.setup)
}
