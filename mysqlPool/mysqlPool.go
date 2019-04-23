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
func NewSQLPool(size int, config MysqlConfig) *MysqlPool {
	hp := new(MysqlPool)
	hp.res = make(chan *sql.DB, size);
	dataSourceName = config.User + ":" + config.Password + "@" + config.Url + "/" + config.Database + "?charset=" + config.Charset
	return hp;
}

//从池子中得倒一个资源
func (p *MysqlPool) GetResource() (conn *sql.DB, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		//fmt.Println("连接池资源" + time.Now().String())
		return r, nil
	default:
		//fmt.Println("新生成资源" + time.Now().String())
		return p.factory()
	}
}

//生成一个资源
func (p *MysqlPool) factory() (conn *sql.DB, err error) {
	client, err := sql.Open("mysql", dataSourceName)
	return client, err
}

//释放资源
func (p *MysqlPool) Release(c *sql.DB) {
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
