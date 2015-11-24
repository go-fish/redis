package redis

import (
	"fmt"
	"strconv"
	"bytes"
	"errors"
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
	ErrTransactionBegin = errors.New(`Fail To Begin Transaction`)
	ErrTransactionAdd	= errors.New(`Fail To Add Command To Transcation Queued`)
	ErrWatchKey			= errors.New(`Fail To Watch Key`)
	ErrUnWatchKey		= errors.New(`Fail To UnWatch Key`)
	ErrDiscard			= errors.New(`Fail To Discard Transaction`)
	ErrExec				= errors.New(`Fail To Exec Transaction`)

	// ErrClosed is the error resulting if the redisPool is closed via RedisPool.Close().
	ErrClosed = errors.New("redisPool is closed")

	endOfLine = []byte{'\r', '\n'}
)

type RedisError error

func intToBytes(num int) []byte {
	var buff [32]byte
	var i = 31
	for {
		buff[i] = byte('0' + num%10)
		i--
		num = num / 10

		if num == 0 {
			break
		}
	}

	return buff[i+1:]
}

func bytesToInt(line []byte) (num int, err error) {
	if len(line) == 0 {
		return -1, errors.New("Malformed Length")
	}
	
	if line[0] == '-' && line[1] == '1' && len(line) == 2 {
		return -1, nil
	}
	
	for _, v := range line {
		num *= 10
		if v < '0' || v > '9' {
			err = fmt.Errorf("illegal line in int")
			num = -1
			return
		}

		num += int(v - '0')
	}

	return
}

func interfaceToBytes(value interface{}) []byte {
	switch value := value.(type) {
	case string:
		return []byte(value)
	case []byte:
		return value
	case int:
		return intToBytes(value)
	case int64:
		return intToBytes(int(value))
	case float64:
		var buff = make([]byte, 16)
		return strconv.AppendFloat(buff, value, 'g', -1, 64)
	case bool:
		if value {
			return []byte{'1'}
		} else {
			return []byte{'0'}
		}
	case nil:
		return nil
	default:
		var buf bytes.Buffer
		fmt.Fprint(&buf, value)
		return buf.Bytes()
	}

	return nil
}

func floatToBytes(value float64) []byte {
	var buff = make([]byte, 16)
	return strconv.AppendFloat(buff, value, 'g', -1, 64)
}

func encodeCommand(command [][]byte) (buff []byte) {
	buff = make([]byte, 0)

	buff = append(buff, '*')
	buff = append(buff, intToBytes(len(command))...)
	buff = append(buff, endOfLine...)

	for _, line := range command {
		var bs = make([]byte, 0)

		bs = append(bs, '$')
		bs = append(bs, intToBytes(len(line))...)
		bs = append(bs, endOfLine...)
		bs = append(bs, line...)
		bs = append(bs, endOfLine...)

		buff = append(buff, bs...)
	}

	return
}

func errorf(line []byte) RedisError {
	return errors.New(string(line))
}