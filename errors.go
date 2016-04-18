package redis

import "errors"

var (
	ErrTCPConn         = errors.New("unspported connection without tcp")
	ErrPoolClosed      = errors.New("connection pool was closed")
	ErrMalformedInt    = errors.New("malformed int")
	ErrMalformedByte   = errors.New("malformed byte")
	ErrUnknownResult   = errors.New("unknown result")
	ErrMalformedLength = errors.New("malformed length")
	ErrNilKey          = errors.New("nil key")
	ErrTooManyArgvs    = errors.New("too many argvs")
	ErrNil             = errors.New("nil returned")
	ErrMalformedArgvs  = errors.New("malformed argvs")
)

type RedisError struct {
	data []byte
}

func (re *RedisError) Error() string {
	return string(re.data)
}
