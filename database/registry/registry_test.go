package registry_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/windyy10/trpc-go-tools/database/http"
	"github.com/windyy10/trpc-go-tools/database/mysql"
	"github.com/windyy10/trpc-go-tools/database/registry"
	"gopkg.in/yaml.v3"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

const (
	yamlFile = `
database:
  mysql:
    providers:
    - name: mysql_test
      service: trpc.mysql.test.db
      timeout: 3000
      target: dsn://{user}:{password}@tcp(ip:port)/{db}?timeout=1s&parseTime=true
      table: {table}
  http:
    providers:
    - name: baidu_page
      service: trpc.http.baidu.page
      timeout: 3000
      url: https://www.baidu.com
`
)

func TestRegistry(t *testing.T) {
	// kafka.RegisterHandleFunc("kafka_handle", func(ctx context.Context, msg *sarama.ConsumerMessage) error {
	//   fmt.Printf("----------kafka consume: key[%s], value[%s]", string(msg.Key), string(msg.Value))
	//   return nil
	// })

	conf := &struct {
		Database yaml.Node `yaml:"database"`
	}{}
	yaml.Unmarshal([]byte(yamlFile), conf)
	err := registry.Setup(&conf.Database)
	require.Nil(t, err)
}

func TestMysql(t *testing.T) {
	db := mysql.Get("mysql_test", context.Background())
	require.NotNil(t, db)
	type DBItem struct {
		Id int `gorm:"column:id"`
	}
	var items []DBItem
	res := db.Limit(1).Find(&items)
	fmt.Print(res.Error)
	fmt.Print(res.RowsAffected)
	fmt.Print(items)
	require.Nil(t, res.Error)
}

// func TestRedis(t *testing.T) {
//   redisClient := redis.Get("redis_test")
//   require.NotNil(t, redisClient)
//   proxy := redis.GetProxy("redis_test")
//   require.NotNil(t, proxy.Client())
//   require.Equal(t, "test:", proxy.PrefixKey())
// }

// func TestKafka(t *testing.T) {
//   kafkaClient := kafka.Get("kafka_test")
//   require.NotNil(t, kafkaClient)
// }

func TestHttpGet(t *testing.T) {
	httpClient := http.Get("baidu_page")
	require.NotNil(t, httpClient)
	m := make(map[string]any)
	rspHeader := &thttp.ClientRspHeader{}
	httpClient.Get(trpc.BackgroundContext(), &m, nil, client.WithRspHead(rspHeader))
	bytes, err := io.ReadAll(rspHeader.Response.Body)
	require.Nil(t, err)
	require.NotEmpty(t, bytes)
}

func TestHttpPostEncoded(t *testing.T) {
	httpClient := http.Get("baidu_page")
	require.NotNil(t, httpClient)
	reqBody := make(map[string]any)
	reqBody["a"] = 1
	data := url.Values{}
	data.Set("key1", "value1")
	data.Set("key2", "value2")
	encodedData := data.Encode()
	m := make(map[string]any)
	err := httpClient.Post(trpc.BackgroundContext(), reqBody, &m, &http.ReqHeader{
		Header: map[string]string{
			"x-abc":        "ddd",
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Path:  "/a/b/c",
		Query: "?a=1&b=2",
		Body:  bytes.NewBufferString(encodedData),
	})
	require.NotNil(t, err)
}

// func TestCos(t *testing.T) {
//   cosClient := fcos.GetProxy("cos_test")
//   require.NotNil(t, cosClient)
// }
