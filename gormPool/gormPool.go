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
	hp.res = make(chan *gorm.DB, size)
	driverName = config.DriverName
	dataSourceName = config.DataSourceName

	for i := 0; i < size; i++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}

	return hp
}

//从池子中得倒一个资源
func (p *GormPool) GetResource() (conn *gorm.DB, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		return r, nil
	}
}

//生成一个资源
func (p *GormPool) factory() (*gorm.DB, error) {
	return gorm.Open(driverName, dataSourceName)
}

//释放资源
func (p *GormPool) Release(c *gorm.DB) {
	if p.close {
		return
	}

	select {
	default:
		p.res <- c
	}
}

// 关闭池子内的连接
func (p *GormPool) Close() {
	p.Lock()
	defer p.Unlock()

	for {
		select {
		case res, ok := <-p.res:
			if ok {
				continue
			}
			defer res.Close()
		default:
			goto End
		}
	}
End:
	p.close = true
}
