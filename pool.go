package redis

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	address           string
	maxConnection     int32
	maxIdleConnection int32
	openNum           int32
	maxIdleTimeout    time.Duration
	idleConns         chan *Conn
	buff              *buffer
	cond              *sync.Cond
	state             int32
}

func Open(address string, maxConnection, maxIdleConnection int, maxIdleTimeout time.Duration) (*Pool, error) {
	p := &Pool{
		address:           address,
		maxConnection:     int32(maxConnection),
		maxIdleConnection: int32(maxIdleConnection),
		maxIdleTimeout:    maxIdleTimeout,
		idleConns:         make(chan *Conn, maxConnection),
		buff:              newBuffer(maxIdleConnection, maxConnection),
		state:             1,
		cond:              sync.NewCond(new(sync.Mutex)),
	}

	//init pool
	for i := 0; i < maxIdleConnection; i++ {
		c, err := p.newConnection()
		if err != nil {
			return nil, err
		}

		p.PutConnection(c)
	}

	p.openNum = int32(maxIdleConnection)
	return p, nil
}

//new connection
func (p *Pool) newConnection() (*Conn, error) {
	c := new(Conn)

	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		return nil, err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, ErrTCPConn
	}

	if err = tcpConn.SetKeepAlive(true); err != nil {
		return nil, err
	}

	c.pool = p
	c.buff = p.buff.get()
	c.rd = bufio.NewReaderSize(conn, 1024)
	c.wd = conn

	return c, nil
}

//get
func (p *Pool) GetConnection() (*Conn, error) {
	for {
		if atomic.LoadInt32(&p.state) == 0 {
			return nil, ErrPoolClosed
		}

		now := time.Now()
		for 0 < len(p.idleConns) {
			conn := <-p.idleConns
			if conn != nil && conn.activeTime.Add(p.maxIdleTimeout).After(now) {
				return conn, nil
			}

			//close timeout connection
			conn.Close()
		}

		//create new connection
		if p.openNum <= p.maxConnection {
			conn, err := p.newConnection()
			if err != nil {
				return nil, err
			}

			atomic.AddInt32(&p.openNum, 1)
			return conn, nil
		}

		//wait free connection
		p.cond.L.Lock()
		for 0 == len(p.idleConns) || p.openNum > p.maxConnection {
			p.cond.Wait()
		}
		p.cond.L.Unlock()
	}
}

//put
func (p *Pool) PutConnection(conn *Conn) {
	if atomic.LoadInt32(&p.state) == 0 {
		return
	}

	if p.idleConns == nil {
		conn.Close()
		return
	}

	//put connection to channel
	conn.activeTime = time.Now()
	select {
	case p.idleConns <- conn:
		p.cond.Signal()
		return

	default:
		conn.Close()
		return
	}
}

//close
func (p *Pool) Close() {
	//change pool's state
	atomic.StoreInt32(&p.state, 0)

	//close all connections
	for conn := range p.idleConns {
		conn.Close()
	}
	close(p.idleConns)

	//free all wait connection
	p.cond.Broadcast()

	//release buffer
	p.buff = nil
}
