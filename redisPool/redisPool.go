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
	hp.res = make(chan redis.Conn, size);

	netWork = config.NetWork
	address = config.Address


	for i := 0; i < size; i ++ {
		conn, _ := hp.factory()
		hp.res <- conn
	}


	return hp;
}

//从池子中得倒一个资源
func (p *RedisPool) GetResource() (conn redis.Conn, err error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, errors.New("pool is close")
		}
		//fmt.Println("连接池资源" + time.Now().String())
		return r, nil
	//default:
		//fmt.Println("新生成资源" + time.Now().String())
		//return p.factory()
	}
}

//生成一个资源
func (p *RedisPool) factory() (conn redis.Conn, err error) {
	client, err := redis.Dial(netWork, address)
	return client, err
}

//释放资源
func (p *RedisPool) Release(c redis.Conn) {
	//////忘了加锁    因为close是线程不安全的
	p.Lock()
	defer p.Unlock()

	if p.close {
		return
	}

	select {
	default:
		p.res <- c
		//fmt.Println("放回连接池资源" + time.Now().String())
		///////这里忘了释放资源的操作了
	}
}
