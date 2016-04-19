package redis

import (
	"bufio"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	pool *Pool
	rd   *bufio.Reader
	wd   net.Conn
	buff *item

	//pipeline
	pipeline []byte
	count    int

	activeTime time.Time
}

//create new connection
func NewConnection(address string) (*Conn, error) {
	c := new(Conn)

	conn, err := net.Dial("tcp", address)
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

	c.buff = newItem()
	c.rd = bufio.NewReaderSize(conn, 1024)
	c.wd = conn

	return c, nil
}

//close connection
func (c *Conn) Close() {
	c.wd.Close()
	c.rd = nil

	//give back buff
	if c.pool != nil {
		c.pool.buff.put(c.buff)

		//release opennum && send signal to cond
		atomic.AddInt32(&c.pool.openNum, -1)
		c.pool.cond.Signal()
	}
}

//exec
func (c *Conn) exec(data []byte) (interface{}, error) {
	//write command
	if _, err := c.wd.Write(data); err != nil {
		//reset buff
		c.buff.reset()
		if c.pipeline != nil {
			c.resetPipeline()
		}

		return nil, err
	}

	//reset buff
	c.buff.reset()
	if c.pipeline != nil {
		c.resetPipeline()
	}

	//read result
	return c.readResult()
}

//read result
func (c *Conn) readResult() (interface{}, error) {
	//reset buff
	defer c.buff.reset()

	data, _, err := c.rd.ReadLine()
	if err != nil {
		return nil, err
	}

	switch data[0] {
	case '+':

		switch {
		case len(data) == 3 && data[1] == 'O' && data[2] == 'K':
			return true, nil

		case len(data) == 5 && data[1] == 'P' && data[2] == 'O' && data[3] == 'N' && data[4] == 'G':
			return true, nil

		default:
			return data[1:], nil
		}

	case '-':
		return nil, &RedisError{data: data[1:]}

	case ':':
		return c.bytesToInt(data[1:])

	case '$':
		count, err := c.bytesToInt(data[1:])
		if err != nil || count < 0 {
			return nil, err
		}

		if count == 0 {
			return nil, nil
		}

		return c.readBulk(count)

	case '*':
		count, err := c.bytesToInt(data[1:])
		if err != nil || count < 0 {
			return nil, err
		}

		if count == 0 {
			return nil, nil
		}

		return c.readMultiBulk(count)

	default:
		return nil, ErrUnknownResult
	}
}

//read bulk
func (c *Conn) readBulk(count int) ([]byte, error) {
	data := c.buff.next(count)

	if _, err := io.ReadFull(c.rd, data); err != nil {
		return nil, err
	}

	//skip \r\n
	if _, _, err := c.rd.ReadLine(); err != nil {
		return nil, err
	}

	return data, nil
}

//read multi bulk
func (c *Conn) readMultiBulk(count int) ([][]byte, error) {
	res := make([][]byte, count)

	for i := 0; i < count; i++ {
		data, _, err := c.rd.ReadLine()
		if err != nil {
			return nil, err
		}

		size, err := c.bytesToInt(data[1:])
		if err != nil {
			return nil, err
		}

		if size < 0 {
			res[i] = nil
			continue
		}

		if res[i], err = c.readBulk(size); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (c *Conn) resetPipeline() {
	c.pipeline = c.pipeline[:0]
	c.count = 0
}
