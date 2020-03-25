package esPool

import (
	"errors"
	"github.com/olivere/elastic"
	"sync"
)

type EsPool struct {
	res chan *elastic.Client
	sync.Mutex
	close bool
}

type EsConfig struct {
	url string // "root:123@/test?charset=utf8"
}

var (
	url string
)

//创建一个pool
func NewEsPool(size int, u string) *EsPool {
	hp := new(EsPool)
	hp.res = make(chan *elastic.Client, size)
	url = u

	for i := 0; i < size; i++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}
	return hp
}

//从池子中得倒一个资源
func (p *EsPool) GetResource() (conn *elastic.Client, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		//fmt.Println("连接池资源" + time.Now().String())
		return r, nil
	}
}

//生成一个资源
func (p *EsPool) factory() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url),
	)
	return client, err
}

//释放资源
func (p *EsPool) Release(c *elastic.Client) {
	if p.close {
		return
	}

	select {
	default:
		p.res <- c
		//fmt.Println("放回连接池资源" + time.Now().String())
	}
}
