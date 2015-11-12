package redis

import (
	"bufio"
	"errors"
	"net"
	"io"
	"syscall"
	"time"
	"bytes"
)

var (
	ErrUndefinedPool    = errors.New(`Undefined Pool`)
	ErrRedisAddr        = errors.New(`Invalid Redis Host Or Port`)
	ErrBadPacket        = errors.New(`Bad Redis Reply Packet`)
	ErrInvalidMaxActive = errors.New(`Invalid maxActive`)
	ErrInvalidAddr      = errors.New(`Invalid Redis Host Or Port`)
	ErrNotConnected     = errors.New(`Client is not connected.`)
	ErrInvalidValue     = errors.New(`Invalid Redis Value`)
	ErrNil              = errors.New(`Nil Returned`)
	ErrGetHeader        = errors.New(`Get Redis Header Error`)
	ErrUnknowRedisReply = errors.New(`Unknow Redis Reply Type`)
	ErrReplyPecket      = errors.New(`Redis Reply Packet Is Nil`)
	ErrTimeout          = errors.New(`Redis Connect Timeout`)

	// ErrClosed is the error resulting if the redisPool is closed via RedisPool.Close().
	ErrClosed = errors.New("redisPool is closed")

	endOfLine = []byte{'\r', '\n'}
)

type Conn struct {
	conn net.Conn

	timeout time.Duration
	timer   *time.Timer

	bw *bufio.Writer
	br *bufio.Reader

	inUse     bool
	isBadConn bool
}

func (this *Conn) close() (err error) {
	err = this.conn.Close()
	this = nil

	return
}

func (this *Conn) redisCommand(command [][]byte) (res interface{}, err error) {
	err = this.sendCommand(command)
	if err != nil {
		return
	}

	return this.decodeCommand()
}

func (this *Conn) sendCommand(command [][]byte) (err error) {
	var buff = encodeCommand(command)

	_, err = this.bw.Write(buff)
	if err != nil {
		return
	}

	return this.flush()
}

func (this *Conn) flush() (err error) {
	err = this.bw.Flush()
	if err == syscall.EPIPE || err == io.EOF {
		this.isBadConn = true
	}

	return
}

func (this *Conn) decodeCommand() (res interface{}, err error) {
	var line []byte

	line, err = this.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '+':
		res = line[1:]
	case '-':
		res = line[1:]
	case ':':
		res = line[1:]
	case '$':
		res, err = this.readBulkData(line[1:])
	case '*':
		res, err = this.readMultiBulkData(line[1:])
	}

	return
}

func (this *Conn) readLine() (line []byte, err error) {
	line, err = this.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	
	if len(line) > 1 && line[len(line) - 2] == '\r' {
		return line, nil
	}
	
	return nil, ErrBadPacket
}

func (this *Conn) getCount(line []byte) (num int, err error) {
	var end = bytes.IndexByte(line, '\r')
	
	return bytesToInt(line[:end])
}

func (this *Conn) readBulkData(line []byte) (res []byte, err error) {
	var num int
	num, err = this.getCount(line)
	if err != nil {
		return nil, err
	}
	
	if num == -1 {
		return line, nil
	}
	
	var n int
	res = make([]byte, num+2)
	n, err = this.br.Read(res)
	if err != nil {
		return nil, err
	}
	
	if n < num {
		return res, nil
	}
	
	return res[:num], nil
}

func (this *Conn) readMultiBulkData(line []byte) (res[][]byte, err error){
	var num int
	num, err = this.getCount(line)
	if err != nil {
		return nil, err
	}
	
	res = make([][]byte, 0, num)
	for i := 0; i < num; i++ {
		buff, err := this.decodeCommand()
		if err != nil {
			return nil, err
		}
		
		if b, ok := buff.([]byte); ok {
			res = append(res, b)
		}
		
		if b1, ok := buff.([][]byte); ok {
			res = append(res, b1...)	
		}
	}
	
	return
}