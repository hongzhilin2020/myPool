package gormPool

import (
	"errors"
	"github.com/jinzhu/gorm"
	"sync"
)

/**
请求链接池
 */
type GormPool struct {
	res chan *gorm.DB
	sync.Mutex
	close bool
}

type GormConfig struct {
	DriverName     string
	DataSourceName string // "root:123@/test?charset=utf8"
}

var (
	driverName,
	dataSourceName string
)

//创建一个pool
func NewGormPool(size int, config GormConfig) *GormPool {
	hp := new(GormPool)
	hp.res = make(chan *gorm.DB, size);
	driverName = config.DriverName
	dataSourceName = config.DataSourceName

	for i := 0; i < size; i ++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}

	return hp;
}

//从池子中得倒一个资源
func (p *GormPool) GetResource() (conn *gorm.DB, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		//fmt.Println("连接池资源" + time.Now().String())
		return r, nil
		//default:
		//	fmt.Println("新生成资源" + time.Now().String())
		//	return p.factory()
	}
}

//生成一个资源
func (p *GormPool) factory() (*gorm.DB, error) {
	engine, err := gorm.Open(driverName, dataSourceName)
	return engine, err
}

//释放资源
func (p *GormPool) Release(c *gorm.DB) {
	p.Lock()
	defer p.Unlock()

	if p.close {
		return
	}

	select {
	default:
		p.res <- c
		//fmt.Println("放回连接池资源" + time.Now().String())
	}
}
