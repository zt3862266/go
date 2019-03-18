package genericpool

import (
	"time"
	"sync"
	"errors"

	"github.com/zt3862266/go/log"
)

var(
	ErrWrongParam = errors.New("pool config invalid")
	ErrConnFailed = errors.New("Connection failed")
	ErrChanClosed = errors.New("channel is closed")
	ErrIdleTimeOut = errors.New("exceed max idle time")
	ErrNilConn = errors.New("connection is nil,reject")
)

// this is connection pool struct
type Pool struct{
	InitSize int
	MaxSize int
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
		return nil,ErrWrongParam
	}
	pool = &Pool{
		InitSize:initSize,
		MaxSize:maxSize,
		MaxIdleTime:maxIdleTime,
		lock:  &sync.Mutex{},
		PoolChan:make(chan connWrapper,maxSize),
		Factory:factory,
	}
	for i:=0;i<initSize;i++{
		conn ,err := factory()
		if err !=nil{
			return nil,ErrConnFailed
		}
		pool.PoolChan <-  connWrapper{
			Conn:conn,
			CreateTime:time.Now(),
		}
	}
	go func(pool *Pool){

		for{
			log.Info("pool InitSize:%d,MaxSize:%d,poolLen:%d",pool.InitSize,pool.MaxSize,len(pool.PoolChan))
			time.Sleep(time.Second *5)
		}

	}(pool)
	return pool,nil
}


func (p *Pool) Get() (conn Conn, err error){

	f  := func() {
		p.lock.Lock()
		defer p.lock.Unlock()
		if len(p.PoolChan) == 0  {
			conn, err := p.Factory()
			if err != nil {
				return
			}
			p.PoolChan <- connWrapper{
				Conn:conn,
				CreateTime:time.Now(),
			}
		}
	}
	for{
		select{
			case connWrapper,ok := <- p.PoolChan:
				if ok {
					if connWrapper.CreateTime.Add(p.MaxIdleTime).Before(time.Now()) {
						connWrapper.Conn.Close()
						log.Info("conn idle timeout:%v",ErrIdleTimeOut)
						continue
					}else {
						return connWrapper.Conn, nil
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
	if conn ==nil{
		return ErrNilConn
	}

	p.lock.Lock()
	select{
		case p.PoolChan <- connWrapper{Conn:conn, CreateTime:time.Now(),}:
			p.lock.Unlock()
			return nil
		default:
			p.lock.Unlock()
			log.Info("conn chan full,close")
			return conn.Close()

	}
	return nil
}

func (p *Pool) close(){
	for connWrapper := range p.PoolChan{
		connWrapper.Conn.Close()
	}
	p.Factory = nil
}
