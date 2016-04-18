package redis

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	True  = []byte{'1'}
	False = []byte{'0'}
)

//interface to bytes
func (c *Conn) interfaceToBytes(val interface{}) []byte {
	switch v := val.(type) {
	case string:
		return []byte(v)

	case []byte:
		return v

	case int:
		return c.intToBytes(v)

	case int8:
		return c.intToBytes(int(v))

	case int16:
		return c.intToBytes(int(v))

	case int32:
		return c.intToBytes(int(v))

	case int64:
		return c.intToBytes(int(v))

	case uint:
		return c.uintToBytes(v)

	case uint8:
		return c.uintToBytes(uint(v))

	case uint16:
		return c.uintToBytes(uint(v))

	case uint32:
		return c.uintToBytes(uint(v))

	case uint64:
		return c.uintToBytes(uint(v))

	case float32:
		return c.float64ToBytes(float64(v))

	case float64:
		return c.float64ToBytes(v)

	case bool:
		if v {
			return True
		}

		return False

	case nil:
		return nil

	default:
		var buf bytes.Buffer
		fmt.Fprint(&buf, v)
		return buf.Bytes()
	}

}

//int to bytes
func (c *Conn) intToBytes(num int) []byte {
	var negative bool
	data := c.buff.next(32)
	i := 31

	if num < 0 {
		negative = true
		num = (num ^ num>>31) - num>>31
	}

	for {
		data[i] = byte('0' + num%10)
		i--
		num = num / 10

		if num == 0 {
			break
		}
	}

	if negative {
		data[i] = '-'
		i--
	}

	return data[i+1:]
}

//uint to bytes
func (c *Conn) uintToBytes(num uint) []byte {
	data := c.buff.next(32)
	i := 31
	for {
		data[i] = byte('0' + num%10)
		i--
		num = num / 10

		if num == 0 {
			break
		}
	}

	return data[i+1:]
}

//int64 to bytes
func (c *Conn) int64ToBytes(num int64) []byte {
	var negative bool
	data := c.buff.next(32)
	i := 31

	if num < 0 {
		negative = true
		num = (num ^ num>>31) - num>>31
	}

	for {
		data[i] = byte('0' + num%10)
		i--
		num = num / 10

		if num == 0 {
			break
		}
	}

	if negative {
		data[i] = '-'
		i--
	}

	return data[i+1:]
}

//bytes to int
func (c *Conn) bytesToInt(data []byte) (int, error) {
	if len(data) == 0 {
		return -1, ErrMalformedInt
	}

	num := 0
	index := 0
	if data[0] == '-' {
		index = 1
	}

	for i := index; i < len(data); i++ {
		num *= 10
		if data[i] < '0' && data[i] > '9' {
			return -1, ErrMalformedByte
		}

		num += int(data[i] - '0')
	}

	if index == 1 {
		return -num, nil
	}

	return num, nil
}

//float64 to bytes
func (c *Conn) float64ToBytes(value float64) []byte {
	data := c.buff.next(32)[:0]
	return strconv.AppendFloat(data, value, 'g', -1, 64)
}

//check buffer
var growPercent float64 = 0.95

func (c *Conn) checkBuffer(offset int, data []byte) []byte {
	if float64(offset)/float64(len(data)) >= growPercent {
		newData := make([]byte, 2*len(data))
		copy(newData, data)
		return newData
	}

	return data
}
