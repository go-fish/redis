package redis

import (
	"fmt"
	"os"
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

var pool *Pool
var pool1 *redigo.Pool

func init() {
	var err error
	if pool, err = Open("127.0.0.1:6379", 10, 5, time.Duration(10)*time.Second); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fn := func() (redigo.Conn, error) {
		return redigo.Dial("tcp4", "127.0.0.1:6379")
	}
	pool1 = redigo.NewPool(fn, 5)
}

func TestDo(t *testing.T) {
	conn, err := pool.GetConnection()
	if err != nil {
		t.Fatalf("wanted nil, get %v", err)
	}

	//do
	val, err := conn.Do("SET", "test111", 111)
	if err != nil {
		t.Fatalf("wanted nil, get %v", err)
	}

	v, ok := val.(bool)
	if !ok {
		t.Fatalf("wanted []byte, get %T", val)
	}

	if !v {
		t.Fatalf("wanted true, get %v", v)
	}

	val, err = conn.Do("GET", "test111")
	if err != nil {
		t.Fatalf("wanted nil, get %v", err)
	}

	v1, ok := val.([]byte)
	if !ok {
		t.Fatalf("wanted int, get %T", val)
	}

	if string(v1) != "111" {
		t.Fatalf("wanted 111, get %v", v1)
	}
}

func Test_LRange10(t *testing.T) {
	conn, err := pool.GetConnection()
	if err != nil {
		t.Fatalf("wanted nil, get %v", err)
	}

	b, err := conn.Do("LRANGE", "hello", 0, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(b)

	s, err := conn.Lrange("hello", 0, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(s)
}

func Test_LRange100(t *testing.T) {
	conn, err := pool.GetConnection()
	if err != nil {
		t.Fatalf("wanted nil, get %v", err)
	}

	s, err := conn.Lrange("hello", 0, 100)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(s)
}
