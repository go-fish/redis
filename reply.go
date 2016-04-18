package redis

import (
	"fmt"
	"strconv"
)

//interface to int
func Int(val interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	switch val := val.(type) {
	case int:
		return val, nil

	case []byte:
		n, err := strconv.ParseInt(string(val), 10, 0)
		return int(n), err

	case nil:
		return -1, nil

	default:
		return 0, fmt.Errorf("unexpected type for Int, got type %T", val)
	}
}

//interface to string
func String(val interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch val := val.(type) {
	case int:
		return strconv.Itoa(val), nil

	case []byte:
		return string(val), nil

	case nil:
		return "", nil

	default:
		return "", fmt.Errorf("unexpected type for String, got type %T", val)
	}
}

//interface to bool
func Bool(val interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	switch val := val.(type) {
	case bool:
		return val, nil

	case int:
		return val == 1, nil

	default:
		return false, fmt.Errorf("unexpected type for Bool, got type %T", val)
	}
}

//interface to int64
func Int64(val interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	switch val := val.(type) {
	case int:
		return int64(val), nil

	case []byte:
		n, err := strconv.ParseInt(string(val), 10, 0)
		return n, err

	case nil:
		return -1, nil

	default:
		return 0, fmt.Errorf("unexpected type for Int64, got type %T", val)
	}
}

//interface to bytes
func Bytes(val interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	switch val := val.(type) {
	case []byte:
		b := make([]byte, len(val))
		copy(b, val)
		return b, nil

	case int:
		return []byte(strconv.Itoa(val)), nil

	case nil:
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected type for []byte, got type %T", val)
	}
}

//interface to []string
func Strings(val interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	switch val := val.(type) {
	case [][]byte:
		strs := make([]string, len(val))
		for i := 0; i < len(val); i++ {
			strs[i] = string(val[i])
		}

		return strs, nil

	case []byte:
		return []string{string(val)}, nil

	case nil:
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected type for []string, got type %T", val)
	}
}

//interface to float64
func Float64(val interface{}, err error) (float64, error) {
	if err != nil {
		return 0.0, err
	}

	switch val := val.(type) {
	case []byte:
		return strconv.ParseFloat(string(val), 64)

	case string:
		return strconv.ParseFloat(val, 64)

	case int:
		return float64(val), nil

	case int64:
		return float64(val), nil

	default:
		return 0.0, fmt.Errorf("unexpected type for float64, got type %T", val)
	}
}

//interface to [][]byte
func BytesArray(val interface{}, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}

	switch val := val.(type) {
	case [][]byte:
		res := make([][]byte, len(val))
		for i := 0; i < len(res); i++ {
			res[i] = make([]byte, len(val[i]))
			copy(res[i], val[i])
		}

		return res, nil

	case []byte:
		res := make([]byte, len(val))
		copy(res, val)

		return [][]byte{res}, nil

	case nil:
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected type for [][]byte, got type %T", val)
	}
}
