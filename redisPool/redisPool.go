package redisPool

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"sync"
)

/**
请求链接池
*/
type RedisPool struct {
	res chan redis.Conn
	sync.Mutex
	close bool
}

type RedisConfig struct {
	NetWork string
	Address string
}

var (
	netWork,
	address string
)

//创建一个pool
func NewRedisPool(size int, config RedisConfig) *RedisPool {
	hp := new(RedisPool)
	hp.res = make(chan redis.Conn, size)
	netWork = config.NetWork
	address = config.Address

	for i := 0; i < size; i++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}

	return hp
}

//从池子中得倒一个资源
func (p *RedisPool) GetResource() (conn redis.Conn, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		return r, nil
	}
}

//生成一个资源
func (p *RedisPool) factory() (conn redis.Conn, err error) {
	return redis.Dial(netWork, address)
}

//释放资源
func (p *RedisPool) Release(c redis.Conn) {
	if p.close {
		return
	}

	select {
	default:
		p.res <- c
	}
}

// 关闭池子内的连接
func (p *RedisPool) Close() {
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
