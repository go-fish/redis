package redis

import (
	"fmt"
	"strconv"
	"bytes"
)

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