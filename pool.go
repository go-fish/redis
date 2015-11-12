package redis

import (
	"bufio"
	"net"
	"sync"
)

type Factory func() (net.Conn, error)

type RedisPool struct {
	conn *Conn

	Address string

	numClosed uint64

	mu           *sync.Mutex //protects following fields
	freeConn     []*Conn
	connRequests []chan *Conn
	numOpen      int
	pendingOpens int

	openerCh chan struct{}
	closed   bool

	MaxIdle   int
	MaxActive int

	ReaderBuffer int
	WriterBuffer int

	factory Factory
}

const (
	defaultReaderBuffer = 4096
	defaultWriterBuffer = 4096
)

func NewRedisPool(maxIdle, maxActive int, address string) (redisPool *RedisPool, err error) {
	redisPool = &RedisPool{
		MaxActive: maxActive,
		MaxIdle:   maxIdle,
		Address:   address,

		ReaderBuffer: defaultReaderBuffer,
		WriterBuffer: defaultWriterBuffer,

		mu: &sync.Mutex{},
	}

	redisPool.factory = func() (nc net.Conn, err error) {
		nc, err = net.Dial("tcp", redisPool.Address)
		if err != nil {
			return
		}

		if tc, ok := nc.(*net.TCPConn); ok {
			err = tc.SetKeepAlive(true)
			if err != nil {
				return
			}
		}

		return
	}

	return
}

func (this *RedisPool) SetBufioBuffer(readerBuffer, writerBuffer int) {
	this.ReaderBuffer = readerBuffer
	this.WriterBuffer = writerBuffer
}

func (this *RedisPool) createNewConn() (*Conn, error) {
	conn, err := this.factory()
	if err != nil {
		return nil, err
	}

	return &Conn{
		conn: conn,

		bw: bufio.NewWriterSize(conn, this.WriterBuffer),
		br: bufio.NewReaderSize(conn, this.ReaderBuffer),
	}, nil
}

func (this *RedisPool) Get() (*Conn, error) {
	return this.get()
}

func (this *RedisPool) get() (*Conn, error) {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return nil, ErrClosed
	}

	var numFree = len(this.freeConn)
	if numFree > 0 {
		var conn = this.freeConn[0]
		copy(this.freeConn, this.freeConn[1:])
		this.freeConn = this.freeConn[:numFree-1]
		conn.inUse = true
		this.mu.Unlock()

		return conn, nil
	}

	//wait for free connection
	if this.MaxActive > 0 && this.numOpen >= this.MaxActive {
		var req = make(chan *Conn, 1)
		this.connRequests = append(this.connRequests, req)
		this.mu.Unlock()

		var ret = <-req

		return ret, nil
	}

	this.numOpen++
	this.mu.Unlock()

	conn, err := this.createNewConn()
	if err != nil {
		this.mu.Lock()
		this.numOpen--
		this.mu.Unlock()

		return nil, err
	}

	this.mu.Lock()
	conn.inUse = true
	this.mu.Unlock()

	return conn, nil
}

func (this *RedisPool) Put(conn *Conn, isBadConn bool) {
	this.mu.Lock()
	if !conn.inUse {
		panic("connection returned that was never out")
	}

	conn.inUse = false

	if isBadConn || conn.isBadConn {
		this.maybeOpenNewConnections()
		this.mu.Unlock()
		conn.close()
		return
	}

	var added = this.putConnDBLocked(conn, isBadConn)
	this.mu.Unlock()

	if !added {
		conn.close()
	}
}

func (this *RedisPool) putConnDBLocked(conn *Conn, isBadConn bool) bool {
	if this.MaxActive > 0 && this.numOpen > this.MaxActive {
		return false
	}

	if c := len(this.connRequests); c > 0 {
		var req = this.connRequests[0]
		copy(this.connRequests, this.connRequests[1:])
		this.connRequests = this.connRequests[:c-1]

		if !isBadConn && !conn.isBadConn {
			conn.inUse = true
		}

		req <- conn
		return true
	} else if !isBadConn && !conn.isBadConn && !this.closed && this.maxIdleConnsLocked() > len(this.freeConn) {
		this.freeConn = append(this.freeConn, conn)
		return true
	}

	return false
}

const defaultMaxIdleConns = 2

func (this *RedisPool) maxIdleConnsLocked() int {
	var n = this.MaxIdle
	switch {
	case n == 0:
		return defaultMaxIdleConns
	case n < 0:
		return 0
	default:
		return n
	}
}

// Assumes db.mu is locked.
// If there are connRequests and the connection limit hasn't been reached,
// then tell the connectionOpener to open new connections.
func (this *RedisPool) maybeOpenNewConnections() {
	var numRequest = len(this.connRequests) - this.pendingOpens
	if this.MaxActive > 0 {
		var numCanOpen = this.MaxActive - (this.numOpen + this.pendingOpens)
		if numRequest > numCanOpen {
			numRequest = numCanOpen
		}
	}

	for numRequest > 0 {
		this.pendingOpens++
		numRequest--
		this.openerCh <- struct{}{}
	}
}

func (this *RedisPool) Close() error {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return nil
	}

	close(this.openerCh)
	var err error
	var fns = make([]func() error, 0, len(this.freeConn))
	for _, conn := range this.freeConn {
		fns = append(fns, conn.close)
	}

	this.freeConn = nil
	this.closed = true
	for _, req := range this.connRequests {
		close(req)
	}
	this.mu.Unlock()

	for _, fn := range fns {
		var err1 = fn()
		if err1 != nil {
			err = err1
		}
	}

	return err
}
