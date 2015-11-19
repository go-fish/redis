package redis

import (
	"bufio"
	"sync"
	"bytes"
	"errors"
	"io"
)

type PipeLine struct {
	bw *bufio.Writer
	br *bufio.Reader

	count int

	mu *sync.Mutex
}

func (this *Conn) PipeLine() *PipeLine {
	return &PipeLine{
		bw:   this.bw,
		br:   this.br,
		mu:   &sync.Mutex{},
	}
}

func (this *PipeLine) Add(command string, args ...interface{}) error {
	var comm = make([][]byte, 0)

	comm = append(comm, []byte(command))

	for i := 0; i < len(args); i++ {
		comm = append(comm, interfaceToBytes(args[i]))
	}

	return this.sendCommand(comm)
}

func (this *PipeLine) Len() int {
	return this.count
}

func (this *PipeLine) sendCommand(command [][]byte) (err error) {
	var buff = encodeCommand(command)

	this.mu.Lock()
	
	this.count++
	_, err = this.bw.Write(buff)
	
	this.mu.Unlock()
	return
}

func (this *PipeLine) Flush() (res []interface{}, err error) {
	err = this.bw.Flush()
	if err != nil {
		return
	}
	
	res = make([]interface{}, 0, this.count)
	for i := 0; i < this.count; i++ {
		var info, err = this.decodeCommand()
		if err != nil {
			return nil, err
		}
		
		res = append(res, info)
	}

	this.count = 0

	return
}

func (this *PipeLine) decodeCommand() (res interface{}, err error) {
	var line []byte

	line, err = this.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '+':
		res = line[1:len(line)-2]
	case '-':
		res = line[1:len(line)-2]
	case ':':
		res = line[1:len(line)-2]
	case '$':
		res, err = this.readBulkData(line[1:])
	case '*':
		res, err = this.readMultiBulkData(line[1:])
	}

	return
}

func (this *PipeLine) readLine() (line []byte, err error) {
	line, err = this.br.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("Read Buffer Size Is Too Small")
	}
	
	if err != nil {
		return nil, err
	}
	
	if len(line) > 1 && line[len(line) - 2] == '\r' {
		return line, nil
	}
	
	return nil, ErrBadPacket
}

func (this *PipeLine) getCount(line []byte) (num int, err error) {
	var end = bytes.IndexByte(line, '\r')
	
	return bytesToInt(line[:end])
}

func (this *PipeLine) readBulkData(line []byte) (res []byte, err error) {
	var num int
	num, err = this.getCount(line)
	if err != nil {
		return nil, err
	}
	
	if num == -1 {
		return line, nil
	}
	
	res = make([]byte, num+2)
	_, err = io.ReadFull(this.br, res)
	if err != nil {
		return nil, err
	}
	
	return res[:num], nil
}

func (this *PipeLine) readMultiBulkData(line []byte) (res[][]byte, err error){
	var num int
	num, err = this.getCount(line)
	if err != nil {
		return nil, err
	}
	
	res = make([][]byte, 0, num+2)
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