package mysql

import (
	"time"

	"github.com/windyy10/trpc-go-tools/database/registry"
	"gorm.io/gorm"
	tgorm "trpc.group/trpc-go/trpc-database/gorm"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

var (
	imp *mysqlImp
)

type mysqlImp struct {
	gormMap map[string]*gorm.DB
}

func (m *mysqlImp) get(name string) *gorm.DB {
	if p, ok := m.gormMap[name]; ok {
		return p
	}
	return nil
}

func (m *mysqlImp) setup(decoder registry.Decoder) error {
	conf := &Config{}
	decoder.Decode(conf)
	log.Debugf("setup conf %+v", conf)
	for _, item := range conf.Providers {
		p, err := tgorm.NewClientProxy(item.Service, client.WithTarget(item.Target),
			client.WithTimeout(time.Duration(item.Timeout)*time.Millisecond))
		if err != nil {
			log.Errorf("new gorm clinet fail %s", err.Error())
			return err
		}
		// 允许带表名
		if item.Table != "" {
			p = p.Table(item.Table)
		}
		m.gormMap[item.Name] = p
		log.Infof("db register %+v ok", item)
	}
	return nil
}

func init() {
	imp = &mysqlImp{
		gormMap: make(map[string]*gorm.DB),
	}
	registry.Registry("mysql", imp.setup)
}
