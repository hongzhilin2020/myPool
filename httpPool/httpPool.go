package httpPool

import (
	"errors"
	"net/http"
	"sync"
)

/**
http 请求链接池
*/
type HttpPool struct {
	res chan *http.Client
	sync.Mutex
	close bool
}

//创建一个pool
func NewHttpPool(size int) *HttpPool {
	hp := new(HttpPool)
	hp.res = make(chan *http.Client, size)
	for i := 0; i < size; i++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}

	return hp
}

//从池子中得倒一个资源
func (p *HttpPool) GetResource() (*http.Client, error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		return r, nil
	}
}

//生成一个资源
func (p *HttpPool) factory() (*http.Client, error) {
	client := new(http.Client)

	return client, nil
}

//释放资源
func (p *HttpPool) Release(c *http.Client) {
	if p.close {
		return
	}

	select {
	default:
		p.res <- c
	}
}
