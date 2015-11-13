package redis

//Keys
func (this *Conn) Keys(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("KEYS"),
			[]byte(key),
		},
	)
}

func (this *Conn) Exists(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("EXISTS"),
			[]byte(key),
		},
	)
}

func (this *Conn) Del(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("DEL"),
			[]byte(key),
		},
	)
}

func (this *Conn) Expire(key string, seconds int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("EXPIRE"),
			intToBytes(int(seconds)),
		},
	)
}

func (this *Conn) Ttl(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("TTL"),
			[]byte(key),
		},
	)
}

func (this *Conn) Type(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("TYPE"),
			[]byte(key),
		},
	)
}

func (this *Conn) Rename(key string, newKey string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("RENAME"),
			[]byte(key),
			[]byte(newKey),
		},
	)
}

func (this *Conn) Renamenx(key string, newKey string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("RENAMENX"),
			[]byte(key),
			[]byte(newKey),
		},
	)
}

func (this *Conn) Randomkey() (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("RANDOMKEY"),
		},
	)
}

//Connect
func (this *Conn) Select(index int) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SELECT"),
			intToBytes(index),
		},
	)
}

func (this *Conn) Auth(password string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("AUTH"),
			[]byte(password),
		},
	)
}

func (this *Conn) Echo(msg string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ECHO"),
			[]byte(msg),
		},
	)
}

func (this *Conn) Ping() (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("PING"),
		},
	)
}

func (this *Conn) Quit() (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("QUIT"),
		},
	)
}

//string
func (this *Conn) Get(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("GET"),
			[]byte(key),
		},
	)
}

func (this *Conn) Set(key string, value interface{}) (interface{}, error) {

	return this.redisCommand(
		[][]byte{
			[]byte("SET"),
			[]byte(key),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Setex(key string, seconds int64, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SETEX"),
			[]byte(key),
			intToBytes(int(seconds)),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Setnx(key string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SETNX"),
			[]byte(key),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Strlen(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("STRLEN"),
			[]byte(key),
		},
	)
}

func (this *Conn) Decr(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("DECR"),
			[]byte(key),
		},
	)
}

func (this *Conn) Decrby(key string, decrement int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("DECRBY"),
			[]byte(key),
			intToBytes(int(decrement)),
		},
	)
}

func (this *Conn) Getset(key string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("GETSET"),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Incr(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("INCR"),
			[]byte(key),
		},
	)
}

func (this *Conn) Incrby(key string, increment int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("INCRBY"),
			[]byte(key),
			intToBytes(int(increment)),
		},
	)
}

func (this *Conn) Mget(keys []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("MGET"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	return this.redisCommand(
		command,
	)
}

// keys 's length should eq to values 's length
func (this *Conn) Mset(keys []string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("MSET"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Msetnx(keys []string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("MSETNX"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Psetex(key string, milliseconds int64, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("PSETEX"),
			intToBytes(int(milliseconds)),
			interfaceToBytes(value),
		},
	)
}

//hash
func (this *Conn) Hdel(key string, fields []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("HDEL"), []byte(key))

	for i := 0; i < len(fields); i++ {
		command = append(command, []byte(fields[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Hexists(key string, field string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HEXISTS"),
			[]byte(key),
			[]byte(field),
		},
	)
}

func (this *Conn) Hget(key string, field string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HGET"),
			[]byte(key),
			[]byte(field),
		},
	)
}

func (this *Conn) Hgetall(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HGETALL"),
			[]byte(key),
		},
	)
}

func (this *Conn) Hincrby(key, field string, increment int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HINCRBY"),
			[]byte(key),
			[]byte(field),
			intToBytes(int(increment)),
		},
	)
}

func (this *Conn) Hkeys(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HKEYS"),
			[]byte(key),
		},
	)
}

func (this *Conn) Hlen(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HLEN"),
			[]byte(key),
		},
	)
}

func (this *Conn) Hmget(key string, fields []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("HMGET"), []byte(key))

	for i := 0; i < len(fields); i++ {
		command = append(command, []byte(fields[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Hmset(key string, fields []string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("HMSET"), []byte(key))

	for i := 0; i < len(fields); i++ {
		command = append(command, []byte(fields[i]))
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Hset(key, field string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HSET"),
			[]byte(key),
			[]byte(field),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Hsetnx(key, field string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HSETNX"),
			[]byte(key),
			[]byte(field),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Hvals(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("HVALS"),
			[]byte(key),
		},
	)
}

//list
func (this *Conn) Blpop(keys []string, timeout int64) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("BLPOP"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	command = append(command, intToBytes(int(timeout)))

	return this.redisCommand(command)
}

func (this *Conn) Brpop(keys []string, timeout int64) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("BRPOP"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	command = append(command, intToBytes(int(timeout)))

	return this.redisCommand(command)
}

func (this *Conn) Lindex(key string, index int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LINDEX"),
			[]byte(key),
			intToBytes(int(index)),
		},
	)
}

func (this *Conn) Llen(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LLEN"),
			[]byte(key),
		},
	)
}

func (this *Conn) Lpop(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LPOP"),
			[]byte(key),
		},
	)
}

func (this *Conn) Lrange(key string, start, stop int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LRANGE"),
			[]byte(key),
			intToBytes(int(start)),
			intToBytes(int(stop)),
		},
	)
}

func (this *Conn) Lpush(key string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("LPUSH"), []byte(key))

	for i := 0; i < len(values); i++ {
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Lpushx(key string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LPUSHX"),
			[]byte(key),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Lrem(key string, count int64, value int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LREM"),
			[]byte(key),
			intToBytes(int(count)),
			intToBytes(int(value)),
		},
	)
}

func (this *Conn) Lset(key string, index int64, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LSET"),
			[]byte(key),
			intToBytes(int(index)),
			interfaceToBytes(value),
		},
	)
}

func (this *Conn) Ltrim(key string, start, stop int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("LTRIM"),
			[]byte(key),
			intToBytes(int(start)),
			intToBytes(int(stop)),
		},
	)
}

func (this *Conn) Rpop(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("RPOP"),
			[]byte(key),
		},
	)
}

func (this *Conn) Rpush(key string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("RPUSH"), []byte(key))

	for i := 0; i < len(values); i++ {
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Rpushx(key string, value interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("RPUSHX"),
			[]byte(key),
			interfaceToBytes(value),
		},
	)
}

//set
func (this *Conn) Sadd(key string, values []interface{}) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("SADD"), []byte(key))

	for i := 0; i < len(values); i++ {
		command = append(command, interfaceToBytes(values[i]))
	}

	return this.redisCommand(
		command,
	)
}

func (this *Conn) Scard(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SCARD"),
			[]byte(key),
		},
	)
}

func (this *Conn) Sdiff(keys []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("SDIFF"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	return this.redisCommand(command)
}

func (this *Conn) Sinter(keys []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("SINTER"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	return this.redisCommand(command)
}

func (this *Conn) Sismember(key string, member interface{}) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SISMEMBER"),
			[]byte(key),
			interfaceToBytes(member),
		},
	)
}

func (this *Conn) Smembers(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SMEMBERS"),
			[]byte(key),
		},
	)
}

func (this *Conn) Spop(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SPOP"),
			[]byte(key),
		},
	)
}

func (this *Conn) Srandmember(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("SRANDMEMBER"),
			[]byte(key),
		},
	)
}

func (this *Conn) Srem(key string, member []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("SREM"))

	for i := 0; i < len(member); i++ {
		command = append(command, []byte(member[i]))
	}

	return this.redisCommand(command)
}

func (this *Conn) Sunion(keys []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("SUNION"))

	for i := 0; i < len(keys); i++ {
		command = append(command, []byte(keys[i]))
	}

	return this.redisCommand(command)

}

//sortedset
func (this *Conn) Zadd(key string, scores []float64, members []string) (interface{}, error) {
	var command [][]byte

	command = append(command, []byte("ZADD"), []byte(key))

	for i := 0; i < len(scores); i++ {
		command = append(command, floatToBytes(scores[i]))
		command = append(command, []byte(members[i]))
	}

	return this.redisCommand(command)
}

func (this *Conn) Zcard(key string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ZCARD"),
			[]byte(key),
		},
	)
}

func (this *Conn) Zcount(key string, min int64, max int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ZCOUNT"),
			[]byte(key),
			intToBytes(int(min)),
			intToBytes(int(max)),
		},
	)
}

func (this *Conn) Zincrby(key string, increment int64, member string) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ZINCRBY"),
			[]byte(key),
			intToBytes(int(increment)),
			[]byte(member),
		},
	)
}

func (this *Conn) Zrange(key string, start int64, stop int64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ZRANGE"),
			[]byte(key),
			intToBytes(int(start)),
			intToBytes(int(stop)),
		},
	)
}

func (this *Conn) Zrangebyscore(key string, start float64, stop float64) (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("ZRANGEBYSCORE"),
			[]byte(key),
			floatToBytes(start),
			floatToBytes(stop),
		},
	)
}

//server
func (this *Conn) Flushdb() (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("FLUSHDB"),
		},
	)
}

func (this *Conn) Flushall() (interface{}, error) {
	return this.redisCommand(
		[][]byte{
			[]byte("FLUSHALL"),
		},
	)
}

func (this *Conn) Do(command string, args ...interface{}) (interface{}, error) {
	var comm = make([][]byte, 0)

	comm = append(comm, []byte(command))

	for i := 0; i < len(args); i++ {
		comm = append(comm, interfaceToBytes(args[i]))
	}

	return this.redisCommand(comm)
}
