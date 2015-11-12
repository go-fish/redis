package redis

import (
	"fmt"
	"strconv"
)

func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case nil:
		return 0, ErrNil
	}
	return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", reply)
}

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case nil:
		return "", ErrNil
	}
	return "", fmt.Errorf("redis: unexpected type for String, got type %T", reply)
}

func Strings(reply interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []byte:
		return []string{string(reply)}, nil
	case string:
		return []string{reply}, nil
	case int:
		return []string{strconv.Itoa(reply)}, nil
	case [][]byte:
		var res []string
		for _, v := range reply {
			res = append(res, string(v))
		}

		return res, nil
	case nil:
		return nil, ErrNil
	}

	return nil, fmt.Errorf("redis: unexpected type for []string, got type %T", reply)
}

func Ints(reply interface{}, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}

	switch reply := reply.(type) {
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return []int{int(n)}, err
	case string:
		n, err := strconv.ParseInt(reply, 10, 0)
		return []int{int(n)}, err
	case int:
		return []int{reply}, err
	case [][]byte:
		var res []int
		for _, v := range reply {
			n, err := strconv.ParseInt(string(v), 10, 0)
			if err != nil {
				return nil, err
			}

			res = append(res, int(n))
		}

		return res, nil
	case nil:
		return nil, ErrNil
	}

	return nil, fmt.Errorf("redis: unexpected type for []int, got type %T", reply)
}
