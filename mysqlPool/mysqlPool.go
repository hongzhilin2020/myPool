package mysqlPool

import (
	"database/sql"
	"errors"
	"sync"
)



/**
http 请求链接池
 */
type MysqlPool struct {
	res chan *sql.DB
	sync.Mutex
	close bool
}

//创建一个pool
func NewSQLPool(size int) *MysqlPool {
	hp := new(MysqlPool)
	hp.res = make(chan *sql.DB,size);
	return hp;
}

//从池子中得倒一个资源
func (p *MysqlPool) GetResource() (conn *sql.DB,err error) {
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
func (p *MysqlPool) factory()  (conn *sql.DB, err error)  {
	//client,err := redis.Dial("tcp", "127.0.0.1:6379")
	client, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/kxxs?charset=utf8mb4")
	return client, err
}

//释放资源
func (p *MysqlPool) Release(c *sql.DB) {
	//////忘了加锁    因为close是线程不安全的
	p.Lock()
	defer p.Unlock()

	if p.close {
		return
	}

	select {
	default:
		p.res<- c
		//fmt.Println("放回连接池资源" + time.Now().String())
		///////这里忘了释放资源的操作了
	}
}

