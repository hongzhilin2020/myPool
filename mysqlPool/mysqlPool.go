package mysqlPool

import (
	"database/sql"
	"errors"
	"sync"
)

/**
请求链接池
*/
type MysqlPool struct {
	res chan *sql.DB
	sync.Mutex
	close bool
}

type MysqlConfig struct {
	User     string
	Password string
	Url      string
	Database string
	Charset  string
}

var dataSourceName string

//创建一个pool
func NewMysqlPool(size int, config MysqlConfig) *MysqlPool {
	hp := new(MysqlPool)
	hp.res = make(chan *sql.DB, size)
	dataSourceName = config.User + ":" + config.Password + "@" + config.Url + "/" + config.Database + "?charset=" + config.Charset
	return hp
}

//从池子中得倒一个资源
func (p *MysqlPool) GetResource() (conn *sql.DB, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		return r, nil
	}
}

//生成一个资源
func (p *MysqlPool) factory() (conn *sql.DB, err error) {
	client, err := sql.Open("mysql", dataSourceName)
	return client, err
}

//释放资源
func (p *MysqlPool) Release(c *sql.DB) {
	if p.close {
		return
	}

	select {
	default:
		p.res <- c
	}
}

// 关闭池子内的连接
func (p *MysqlPool) Close() {
	p.Lock()
	defer p.Unlock()

	for {
		select {
		case res, ok := <-p.res:
			if ok {
				continue
			}
			// 先把标记设置为关闭状态
			defer res.Close()
		default:
			goto End
		}
	}
End:
	p.close = true
}
