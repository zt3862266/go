package genericpool

import (
	"time"
	"sync"
	"github.com/pkg/errors"
	"github.com/zt3862266/go/log"
)

var(
	ErrWrongParam = errors.New("pool config invalid")
	ErrConnFailed = errors.New("Connection failed")
	ErrChanClosed = errors.New("channel is closed")
	ErrIdletimeOut = errors.New("exceed max idle time")
)

// this is connection pool struct
type Pool struct{
	InitSize int
	MaxSize int
	Size int
	MaxIdleTime time.Duration
	Factory func() (Conn,error)
	lock *sync.Mutex
	PoolChan chan connWrapper

}

// conn struct
type Conn interface{
	Close() error
}

// connection的一层wrapper,增加了连接创造时间
type connWrapper struct{
	CreateTime time.Time
	Conn
}


func NewPool(factory func() (Conn,error) ,initSize int,maxSize int,maxIdleTime time.Duration) (pool *Pool,err error){

	if initSize < 1 || maxSize < 1 || initSize > maxSize {
		return pool,ErrWrongParam
	}
	pool = &Pool{
		InitSize:initSize,
		MaxSize:maxSize,
		MaxIdleTime:maxIdleTime,
		lock:  &sync.Mutex{},
		PoolChan:make(chan connWrapper,maxSize),
		Size:0,
		Factory:factory,
	}
	for i:=0;i<initSize;i++{
		conn ,err := factory()
		if err !=nil{
			return pool,ErrConnFailed
		}
		value,ok :=conn.(Conn)

		if ok{
			pool.Size = pool.Size+1
			pool.PoolChan <-  connWrapper{
				Conn:value,
				CreateTime:time.Now(),
			}
		}
	}
	go func(pool *Pool){

		for{
			log.Info("pool InitSize:%d,MaxSize:%d,Size:%d,poolLen:%d",pool.InitSize,pool.MaxSize,pool.Size, len(pool.PoolChan))
			time.Sleep(time.Second *5)
		}

	}(pool)
	return pool,nil
}

func (p *Pool)resize(step int){
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Size = p.Size + step
}

func (p *Pool) Get() (conn Conn, err error){

	f  := func() {
		//chan 为空 而且未达到 MaxSize,可以继续创建 Conn
		if len(p.PoolChan) == 0 && p.Size < p.MaxSize {
			p.lock.Lock()
			defer p.lock.Unlock()

			conn, err := p.Factory()
			if err != nil {
				return
			}
			value, ok := conn.(Conn)
			if ok {
				p.Size = p.Size+1
				p.PoolChan <- connWrapper{
					Conn:value,
					CreateTime:time.Now(),
				}
			}
		}
	}
	for{
		select{
			case connWrapper,ok := <- p.PoolChan:
				if ok {
					if time.Since(connWrapper.CreateTime) > p.MaxIdleTime{
						p.resize(-1)
						connWrapper.Conn.Close()
						log.Info("conn idle timeout:%v",ErrIdletimeOut)
						continue
					}else {
						return conn, nil
					}
				}else{
					return nil,ErrChanClosed
				}
		default:
			f()
		}

	}
}

func (p *Pool) Release(conn Conn) error{
	p.PoolChan <- connWrapper{
		Conn:conn,
		CreateTime:time.Now(),
	}
	return nil
}

func (p *Pool) close(){
	for connWrapper := range p.PoolChan{
		connWrapper.Conn.Close()
	}
	p.Size =0
	p.Factory = nil
}
