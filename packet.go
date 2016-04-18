package redis

//write command
//remember that the result with []byte stop being valid at the next read or write
func (c *Conn) Do(command string, argvs ...interface{}) (interface{}, error) {
	data := c.buff.next(1024)

	size := len(argvs) + 1
	data[0] = '*'
	offset := 1
	offset += copy(data[offset:], c.intToBytes(size))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(command)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], command)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//argvs
	for i := 0; i < len(argvs); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(argvs[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return c.exec(data[:offset])
}

//Ping
func (c *Conn) Ping() (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '1'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'P'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'G'
	data[12] = '\r'
	data[13] = '\n'

	return Bool(c.exec(data[:14]))
}

//set
func (c *Conn) Set(key string, value interface{}) error {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '3'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = '\r'
	data[12] = '\n'
	offset := 13

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	_, err := c.exec(data[:offset])
	return err
}

//Get
func (c *Conn) Get(key string) ([]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '3'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'G'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = '\r'
	data[12] = '\n'
	offset := 13

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bytes(c.exec(data[:offset]))
}

//DEL
func (c *Conn) Del(keys ...string) (int64, error) {
	if keys == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '3'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'D'
	data[offset+5] = 'E'
	data[offset+6] = 'L'
	data[offset+7] = '\r'
	data[offset+8] = '\n'
	offset += 9

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Dump
func (c *Conn) Dump(key string) ([]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'D'
	data[9] = 'U'
	data[10] = 'M'
	data[11] = 'P'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bytes(c.exec(data[:offset]))
}

//Exists
func (c *Conn) Exists(key string) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'E'
	data[9] = 'X'
	data[10] = 'I'
	data[11] = 'S'
	data[12] = 'T'
	data[13] = 'S'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Expire
func (c *Conn) Expire(key string, seconds int) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'E'
	data[9] = 'X'
	data[10] = 'P'
	data[11] = 'I'
	data[12] = 'R'
	data[13] = 'E'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.intToBytes(seconds)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Expireat
func (c *Conn) Expireat(key string, unixtime int64) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'E'
	data[9] = 'X'
	data[10] = 'P'
	data[11] = 'I'
	data[12] = 'R'
	data[13] = 'E'
	data[14] = 'A'
	data[15] = 'T'
	data[16] = '\r'
	data[17] = '\n'
	offset := 18

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(unixtime)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//keys
func (c *Conn) Keys(pattern string) ([]string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'K'
	data[9] = 'E'
	data[10] = 'Y'
	data[11] = 'S'
	data[11] = '\r'
	data[12] = '\n'
	offset := 13

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(pattern)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], pattern)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Strings(c.exec(data[:offset]))
}

//Pexpire
func (c *Conn) Pexpire(key string, milliseconds int64) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'P'
	data[9] = 'E'
	data[10] = 'X'
	data[11] = 'P'
	data[12] = 'I'
	data[13] = 'R'
	data[14] = 'E'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(milliseconds)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Pexpireat
//Pexpire
func (c *Conn) Pexpireat(key string, millisecondsTimestamp int64) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '9'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'P'
	data[9] = 'E'
	data[10] = 'X'
	data[11] = 'P'
	data[12] = 'I'
	data[13] = 'R'
	data[14] = 'E'
	data[15] = 'A'
	data[16] = 'T'
	data[17] = '\r'
	data[18] = '\n'
	offset := 19

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(millisecondsTimestamp)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//PTTL
func (c *Conn) Pttl(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'P'
	data[9] = 'T'
	data[10] = 'T'
	data[11] = 'L'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//TTL
func (c *Conn) Ttl(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '3'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'T'
	data[9] = 'T'
	data[10] = 'L'
	data[11] = '\r'
	data[12] = '\n'
	offset := 13

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Randomkey
func (c *Conn) Randomkey() (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '1'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '9'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'A'
	data[10] = 'N'
	data[11] = 'D'
	data[12] = 'O'
	data[13] = 'M'
	data[14] = 'K'
	data[15] = 'E'
	data[16] = 'Y'
	data[17] = '\r'
	data[18] = '\n'

	return String(c.exec(data[:19]))
}

//Rename
func (c *Conn) Rename(key, newKey string) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'E'
	data[10] = 'N'
	data[11] = 'A'
	data[12] = 'M'
	data[13] = 'E'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//new key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(newKey)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], newKey)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Renamenx
func (c *Conn) Renamenx(key, newKey string) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'E'
	data[10] = 'N'
	data[11] = 'A'
	data[12] = 'M'
	data[13] = 'E'
	data[14] = 'N'
	data[15] = 'X'
	data[16] = '\r'
	data[17] = '\n'
	offset := 18

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//new key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(newKey)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], newKey)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Type
func (c *Conn) Type(key string) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'T'
	data[9] = 'Y'
	data[10] = 'P'
	data[11] = 'E'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

/********** STRING **********/
//Append
func (c *Conn) Append(key, value string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'A'
	data[9] = 'P'
	data[10] = 'P'
	data[11] = 'E'
	data[12] = 'N'
	data[13] = 'D'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//new key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(value)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], value)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Decr
func (c *Conn) Decr(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'D'
	data[9] = 'E'
	data[10] = 'C'
	data[11] = 'R'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Decrby
func (c *Conn) Decrby(key string, value int) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'D'
	data[9] = 'E'
	data[10] = 'C'
	data[11] = 'R'
	data[12] = 'B'
	data[13] = 'Y'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.intToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Getrange
func (c *Conn) Getrange(key string, start int, end int) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'G'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = 'R'
	data[12] = 'A'
	data[13] = 'N'
	data[14] = 'G'
	data[15] = 'E'
	data[16] = '\r'
	data[17] = '\n'
	offset := 18

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.intToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//end
	data[offset] = '$'
	offset++

	val = c.intToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Getset
func (c *Conn) GetSet(key string, value interface{}) (interface{}, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'G'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = 'S'
	data[12] = 'E'
	data[13] = 'T'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return c.exec(data[:offset])
}

//Incr
func (c *Conn) Incr(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'I'
	data[9] = 'N'
	data[10] = 'C'
	data[11] = 'R'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Incrby
func (c *Conn) Incrby(key string, value int64) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'I'
	data[9] = 'N'
	data[10] = 'C'
	data[11] = 'R'
	data[12] = 'B'
	data[13] = 'Y'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Incrbyfloat
func (c *Conn) IncrbyFloat(key string, value float64) (float64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '1'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'C'
	data[12] = 'R'
	data[13] = 'B'
	data[14] = 'Y'
	data[15] = 'F'
	data[16] = 'L'
	data[17] = 'O'
	data[18] = 'A'
	data[19] = 'T'
	data[20] = '\r'
	data[21] = '\n'
	offset := 22

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.float64ToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Float64(c.exec(data[:offset]))
}

//Mget
func (c *Conn) Mget(keys []string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'M'
	data[offset+5] = 'G'
	data[offset+6] = 'E'
	data[offset+7] = 'T'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Mset
func (c *Conn) Mset(keys []string, values []interface{}) error {
	if len(keys) != len(values) {
		return ErrMalformedArgvs
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+len(values)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'M'
	data[offset+5] = 'S'
	data[offset+6] = 'E'
	data[offset+7] = 'T'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//keys && values
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)

		data[offset] = '$'
		offset++
		v := c.interfaceToBytes(values[i])
		offset += copy(data[offset:], c.intToBytes(len(v)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], v)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	_, err := c.exec(data[:offset])
	return err
}

//Msetnx
func (c *Conn) Msetnx(keys []string, values []interface{}) error {
	if len(keys) != len(values) {
		return ErrMalformedArgvs
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+len(values)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '6'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'M'
	data[offset+5] = 'S'
	data[offset+6] = 'E'
	data[offset+7] = 'T'
	data[offset+8] = 'N'
	data[offset+9] = 'X'
	data[offset+10] = '\r'
	data[offset+11] = '\n'
	offset += 12

	//keys && values
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)

		data[offset] = '$'
		offset++
		v := c.interfaceToBytes(values[i])
		offset += copy(data[offset:], c.intToBytes(len(v)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], v)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	_, err := c.exec(data[:offset])
	return err
}

//Psetex
func (c *Conn) Psetex(key string, milliseconds int64, value interface{}) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'P'
	data[9] = 'S'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = 'E'
	data[13] = 'X'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//seconds
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(milliseconds)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val = c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Setex
func (c *Conn) Setex(key string, seconds int64, value interface{}) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = 'E'
	data[12] = 'X'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//seconds
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(seconds)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val = c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Setnx
func (c *Conn) Setnx(key string, value interface{}) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = 'N'
	data[12] = 'X'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Setrange
func (c *Conn) Setrange(key string, offset int64, value interface{}) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'E'
	data[10] = 'T'
	data[11] = 'R'
	data[12] = 'A'
	data[13] = 'N'
	data[14] = 'G'
	data[15] = 'E'
	data[16] = '\r'
	data[17] = '\n'
	pos := 18

	//key
	data[pos] = '$'
	pos++
	pos += copy(data[pos:], c.intToBytes(len(key)))
	data[pos] = '\r'
	data[pos+1] = '\n'
	pos += 2
	pos += copy(data[pos:], key)
	data[pos] = '\r'
	data[pos] = '\n'
	pos += 2

	//offset
	data[pos] = '$'
	pos++

	val := c.int64ToBytes(offset)
	pos += copy(data[pos:], c.intToBytes(len(val)))
	data[pos] = '\r'
	data[pos+1] = '\n'
	pos += 2
	pos += copy(data[pos:], val)
	data[pos] = '\r'
	data[pos+1] = '\n'
	pos += 2

	//value
	data[pos] = '$'
	pos++

	val = c.interfaceToBytes(value)
	pos += copy(data[pos:], c.intToBytes(len(val)))
	data[pos] = '\r'
	data[pos+1] = '\n'
	pos += 2
	pos += copy(data[pos:], val)
	data[pos] = '\r'
	data[pos+1] = '\n'
	pos += 2

	return Int64(c.exec(data[:pos]))
}

//Strlen
func (c *Conn) Strlen(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'T'
	data[10] = 'R'
	data[11] = 'L'
	data[12] = 'E'
	data[13] = 'N'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Scan
func (c *Conn) Scan(cursor int64, arguments ...interface{}) ([][]byte, error) {
	if arguments == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(arguments)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'C'
	data[offset+6] = 'A'
	data[offset+7] = 'N'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//cursor
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(cursor)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//values
	for i := 0; i < len(arguments); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(arguments[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

/******** HASH ********/
func (c *Conn) Hdel(key string, fields ...string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(fields)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'H'
	data[offset+5] = 'D'
	data[offset+6] = 'L'
	data[offset+7] = 'E'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//fields
	for i := 0; i < len(fields); i++ {
		data[offset] = '$'
		offset++

		offset += copy(data[offset:], c.intToBytes(len(fields[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], fields[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//HExists
func (c *Conn) Hexists(key, field string) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'E'
	data[10] = 'X'
	data[11] = 'I'
	data[12] = 'S'
	data[13] = 'T'
	data[14] = 'S'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Hget
func (c *Conn) Hget(key, field string) ([]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'G'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bytes(c.exec(data[:offset]))
}

//hgetall
func (c *Conn) Hgetall(key string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'G'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = 'A'
	data[13] = 'L'
	data[14] = 'L'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return BytesArray(c.exec(data[:offset]))
}

//Hincrby
func (c *Conn) Hincrby(key, field string, increment int64) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'C'
	data[12] = 'R'
	data[13] = 'B'
	data[14] = 'Y'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//increment
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(increment)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Hincrbyfloat
func (c *Conn) Hincrbyfloat(key, field string, increment float64) (float64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '2'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'H'
	data[10] = 'I'
	data[11] = 'N'
	data[12] = 'C'
	data[13] = 'R'
	data[14] = 'B'
	data[15] = 'Y'
	data[16] = 'F'
	data[17] = 'L'
	data[18] = 'O'
	data[19] = 'A'
	data[20] = 'T'
	data[21] = '\r'
	data[22] = '\n'
	offset := 23

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.float64ToBytes(increment)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Float64(c.exec(data[:offset]))
}

//Hkeys
func (c *Conn) Hkeys(key string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'K'
	data[10] = 'E'
	data[11] = 'Y'
	data[12] = 'S'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return BytesArray(c.exec(data[:offset]))
}

//Hlen
func (c *Conn) Hlen(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'L'
	data[10] = 'E'
	data[11] = 'N'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Hmget
func (c *Conn) Hmget(key string, fields ...string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(fields)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'H'
	data[offset+5] = 'M'
	data[offset+6] = 'G'
	data[offset+7] = 'E'
	data[offset+8] = 'T'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//fields
	for i := 0; i < len(fields); i++ {
		data[offset] = '$'
		offset++

		offset += copy(data[offset:], c.intToBytes(len(fields[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], fields[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Hmset
func (c *Conn) Hmset(key string, fields []string, values []interface{}) error {
	if len(fields) != len(values) {
		return ErrMalformedArgvs
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(fields)+len(values)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'H'
	data[offset+5] = 'M'
	data[offset+6] = 'S'
	data[offset+7] = 'E'
	data[offset+8] = 'T'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//fields && values
	for i := 0; i < len(fields); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(fields[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], fields[i])
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)

		data[offset] = '$'
		offset++
		v := c.interfaceToBytes(values[i])
		offset += copy(data[offset:], c.intToBytes(len(v)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], v)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	_, err := c.exec(data[:offset])
	return err
}

//Hset
func (c *Conn) Hset(key, field string, value interface{}) error {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'S'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	_, err := c.exec(data[:offset])
	return err
}

//Hsetnx
func (c *Conn) Hsetnx(key, field string, value interface{}) error {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'S'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = 'N'
	data[13] = 'X'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//field
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(field)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], field)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	_, err := c.exec(data[:offset])
	return err
}

//hvals
func (c *Conn) Hvals(key string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'H'
	data[9] = 'V'
	data[10] = 'A'
	data[11] = 'L'
	data[12] = 'S'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return BytesArray(c.exec(data[:offset]))
}

/******** LIST ********/
//Blpop
func (c *Conn) Blpop(keys []string, timeout int64) ([]string, error) {
	if keys == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'B'
	data[offset+5] = 'L'
	data[offset+6] = 'P'
	data[offset+7] = 'O'
	data[offset+8] = 'P'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	//timeout
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(timeout)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Strings(c.exec(data[:offset]))
}

//Brpop
func (c *Conn) Brpop(keys []string, timeout int64) ([]string, error) {
	if keys == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'B'
	data[offset+5] = 'R'
	data[offset+6] = 'P'
	data[offset+7] = 'O'
	data[offset+8] = 'P'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	//timeout
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(timeout)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Strings(c.exec(data[:offset]))
}

//Brpoplpush
func (c *Conn) Brpoplpush(source, destination string, timeout int64) ([]string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '0'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'B'
	data[10] = 'R'
	data[11] = 'P'
	data[12] = 'O'
	data[13] = 'P'
	data[14] = 'L'
	data[15] = 'P'
	data[16] = 'U'
	data[17] = 'S'
	data[18] = 'H'
	data[19] = '\r'
	data[20] = '\n'
	offset := 21

	//source
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(source)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], source)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//timeout
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(timeout)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Strings(c.exec(data[:offset]))
}

//Lindex
func (c *Conn) Lindex(key string, index int64) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'D'
	data[12] = 'E'
	data[13] = 'X'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//index
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(index)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Linsert
func (c *Conn) Linsert(key, argv, pivot, value string) (int64, error) {
	if argv != "BEFORE" || argv != "AFTER" {
		return -1, ErrMalformedArgvs
	}

	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '5'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'S'
	data[12] = 'E'
	data[13] = 'R'
	data[14] = 'T'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//argv
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(argv)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], argv)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//pivot
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(pivot)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], pivot)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(value)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], value)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Llen
func (c *Conn) Llen(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'L'
	data[10] = 'E'
	data[11] = 'N'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Lpop
func (c *Conn) Lpop(key string) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'P'
	data[10] = 'O'
	data[11] = 'P'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Lpush
func (c *Conn) Lpush(key string, values ...interface{}) (int64, error) {
	if values == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(values)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'L'
	data[offset+5] = 'P'
	data[offset+6] = 'U'
	data[offset+7] = 'S'
	data[offset+8] = 'H'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//values
	for i := 0; i < len(values); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(values[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Lpushx
func (c *Conn) Lpushx(key string, value interface{}) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'P'
	data[10] = 'U'
	data[11] = 'S'
	data[12] = 'H'
	data[13] = 'X'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	b := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(b)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	offset += copy(data[offset:], b)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Lrange
func (c *Conn) Lrange(key string, start, end int64) ([]string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'R'
	data[10] = 'A'
	data[11] = 'N'
	data[12] = 'G'
	data[13] = 'E'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//end
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(end)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Strings(c.exec(data[:offset]))
}

//Lrem
func (c *Conn) Lrem(key string, count uint8, value string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'R'
	data[10] = 'E'
	data[11] = 'M'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//count
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = count
	data[offset+5] = '\r'
	data[offset+6] = '\n'
	offset += 7

	//value
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(value)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], value)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Lset
func (c *Conn) Lset(key string, index int64, value string) error {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'S'
	data[10] = 'E'
	data[11] = 'T'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//index
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(index)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(value)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], value)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	_, err := c.exec(data[:offset])
	return err
}

//Ltrim
func (c *Conn) Ltrim(key string, start, end int64) error {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'L'
	data[9] = 'T'
	data[10] = 'R'
	data[11] = 'I'
	data[12] = 'M'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//end
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(end)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	_, err := c.exec(data[:offset])
	return err
}

//Rpop
func (c *Conn) Rpop(key string) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'P'
	data[10] = 'O'
	data[11] = 'P'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Rpoplpush
func (c *Conn) Rpoplpush(source, destination string) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '9'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'P'
	data[10] = 'O'
	data[11] = 'P'
	data[12] = 'L'
	data[13] = 'P'
	data[14] = 'U'
	data[15] = 'S'
	data[16] = 'H'
	data[17] = '\r'
	data[18] = '\n'
	offset := 19

	//source
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(source)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], source)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Rpush
func (c *Conn) Rpush(key string, values ...interface{}) (int64, error) {
	if values == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(values)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'R'
	data[offset+5] = 'P'
	data[offset+6] = 'U'
	data[offset+7] = 'S'
	data[offset+8] = 'H'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//values
	for i := 0; i < len(values); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(values[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Rpushx
func (c *Conn) Rpushx(key string, value interface{}) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'R'
	data[9] = 'P'
	data[10] = 'U'
	data[11] = 'S'
	data[12] = 'H'
	data[13] = 'X'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//value
	data[offset] = '$'
	offset++

	b := c.interfaceToBytes(value)
	offset += copy(data[offset:], c.intToBytes(len(b)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	offset += copy(data[offset:], b)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

/******** SET ********/
//Sadd
func (c *Conn) Sadd(key string, members ...interface{}) (int64, error) {
	if members == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(members)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'A'
	data[offset+6] = 'D'
	data[offset+7] = 'D'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//members
	for i := 0; i < len(members); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(members[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Scard
func (c *Conn) Scard(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'C'
	data[10] = 'A'
	data[11] = 'R'
	data[12] = 'D'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Sdiff
func (c *Conn) Sdiff(keys ...string) ([][]byte, error) {
	if keys == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'D'
	data[offset+6] = 'I'
	data[offset+7] = 'F'
	data[offset+8] = 'F'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Sdiffstore
func (c *Conn) Sdiffstore(destination string, keys ...string) (int64, error) {
	if keys == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '0'
	data[offset+3] = '\r'
	data[offset+4] = '\n'
	data[offset+5] = 'S'
	data[offset+6] = 'D'
	data[offset+7] = 'I'
	data[offset+8] = 'F'
	data[offset+9] = 'F'
	data[offset+10] = 'S'
	data[offset+11] = 'T'
	data[offset+12] = 'O'
	data[offset+13] = 'R'
	data[offset+14] = 'E'
	data[offset+15] = '\r'
	data[offset+16] = '\n'
	offset += 17

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Sinter
func (c *Conn) Sinter(keys ...string) ([][]byte, error) {
	if keys == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '6'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'I'
	data[offset+6] = 'N'
	data[offset+7] = 'T'
	data[offset+8] = 'E'
	data[offset+9] = 'R'
	data[offset+10] = '\r'
	data[offset+11] = '\n'
	offset += 12

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Sinterstore
func (c *Conn) Sinterstore(destination string, keys ...string) (int64, error) {
	if keys == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '1'
	data[offset+3] = '\r'
	data[offset+4] = '\n'
	data[offset+5] = 'S'
	data[offset+6] = 'I'
	data[offset+7] = 'N'
	data[offset+8] = 'T'
	data[offset+9] = 'E'
	data[offset+10] = 'R'
	data[offset+11] = 'S'
	data[offset+12] = 'T'
	data[offset+13] = 'O'
	data[offset+14] = 'R'
	data[offset+15] = 'E'
	data[offset+16] = '\r'
	data[offset+17] = '\n'
	offset += 18

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Sismember
func (c *Conn) Sismember(key string, member interface{}) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '9'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'I'
	data[10] = 'S'
	data[11] = 'M'
	data[12] = 'E'
	data[13] = 'M'
	data[14] = 'B'
	data[15] = 'E'
	data[16] = 'R'
	data[17] = '\r'
	data[18] = '\n'
	offset := 19

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Smembers
func (c *Conn) Smembers(key string) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'M'
	data[10] = 'E'
	data[11] = 'M'
	data[12] = 'B'
	data[13] = 'E'
	data[14] = 'R'
	data[15] = 'S'
	data[16] = '\r'
	data[17] = '\n'
	offset := 18

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return BytesArray(c.exec(data[:offset]))
}

//Smove
func (c *Conn) Smove(source, destination string, member interface{}) (bool, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'M'
	data[10] = 'O'
	data[11] = 'V'
	data[12] = 'E'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//source
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(source)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], source)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Bool(c.exec(data[:offset]))
}

//Spop
func (c *Conn) Spop(key string) ([]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '4'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'S'
	data[9] = 'P'
	data[10] = 'O'
	data[11] = 'P'
	data[12] = '\r'
	data[13] = '\n'
	offset := 14

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Bytes(c.exec(data[:offset]))
}

//Srandmember
func (c *Conn) Srandmember(key string, count int64) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '1'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'S'
	data[10] = 'R'
	data[11] = 'A'
	data[12] = 'N'
	data[13] = 'D'
	data[14] = 'M'
	data[15] = 'E'
	data[16] = 'M'
	data[17] = 'B'
	data[18] = 'E'
	data[19] = 'R'
	data[20] = '\r'
	data[21] = '\n'
	offset := 22

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//count
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(count)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return BytesArray(c.exec(data[:offset]))
}

//Srem
func (c *Conn) Srem(key string, members ...interface{}) (int64, error) {
	data := c.buff.next(1024)

	size := len(members) + 2
	data[0] = '*'
	offset := 1
	offset += copy(data[offset:], c.intToBytes(size))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'R'
	data[offset+6] = 'E'
	data[offset+7] = 'M'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//members
	for i := 0; i < len(members); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(members[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Sunion
func (c *Conn) Sunion(keys ...string) ([][]byte, error) {
	if keys == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+1)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '6'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'U'
	data[offset+6] = 'N'
	data[offset+7] = 'I'
	data[offset+8] = 'O'
	data[offset+9] = 'N'
	data[offset+10] = '\r'
	data[offset+11] = '\n'
	offset += 12

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Sunionstore
func (c *Conn) Sunionstore(destination string, keys ...string) (int64, error) {
	if keys == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(keys)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '1'
	data[offset+3] = '\r'
	data[offset+4] = '\n'
	data[offset+5] = 'S'
	data[offset+6] = 'U'
	data[offset+7] = 'N'
	data[offset+8] = 'I'
	data[offset+9] = 'O'
	data[offset+10] = 'N'
	data[offset+11] = 'S'
	data[offset+12] = 'T'
	data[offset+13] = 'O'
	data[offset+14] = 'R'
	data[offset+15] = 'E'
	data[offset+16] = '\r'
	data[offset+17] = '\n'
	offset += 18

	//destination
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(destination)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], destination)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//keys
	for i := 0; i < len(keys); i++ {
		data[offset] = '$'
		offset++
		offset += copy(data[offset:], c.intToBytes(len(keys[i])))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2
		offset += copy(data[offset:], keys[i])
		data[offset] = '\r'
		data[offset] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Sscan
func (c *Conn) Sscan(key string, cursor int64, arguments ...interface{}) ([][]byte, error) {
	if arguments == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(arguments)+3)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'S'
	data[offset+5] = 'S'
	data[offset+6] = 'C'
	data[offset+7] = 'A'
	data[offset+8] = 'N'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//cursor
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(cursor)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//values
	for i := 0; i < len(arguments); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(arguments[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

/******** SORTEDSET ********/
//Zadd
func (c *Conn) Zadd(key string, members ...interface{}) (int64, error) {
	if members == nil {
		return -1, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(members)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'Z'
	data[offset+5] = 'A'
	data[offset+6] = 'D'
	data[offset+7] = 'D'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//members
	for i := 0; i < len(members); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(members[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Zcard
func (c *Conn) Zcard(key string) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '2'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'C'
	data[10] = 'A'
	data[11] = 'R'
	data[12] = 'D'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zcount
func (c *Conn) Zcount(key string, min, max int64) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'C'
	data[10] = 'O'
	data[11] = 'U'
	data[12] = 'N'
	data[13] = 'T'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//min
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(min)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//max
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(max)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zincrby
func (c *Conn) Zincrby(key string, increment int64, member interface{}) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '7'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'I'
	data[10] = 'N'
	data[11] = 'C'
	data[12] = 'R'
	data[13] = 'B'
	data[14] = 'Y'
	data[15] = '\r'
	data[16] = '\n'
	offset := 17

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//increment
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(increment)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val = c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Zrange
func (c *Conn) Zrange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'R'
	data[10] = 'A'
	data[11] = 'N'
	data[12] = 'G'
	data[13] = 'E'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//stop
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(stop)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//withScores
	if withScores {
		data[offset] = '$'
		data[offset+1] = '1'
		data[offset+2] = '1'
		data[offset+3] = '\r'
		data[offset+4] = '\n'
		data[offset+5] = 'W'
		data[offset+6] = 'I'
		data[offset+7] = 'T'
		data[offset+8] = 'H'
		data[offset+9] = 'S'
		data[offset+10] = 'C'
		data[offset+11] = 'O'
		data[offset+12] = 'R'
		data[offset+13] = 'E'
		data[offset+14] = 'S'
		data[offset+15] = '\r'
		data[offset+16] = '\n'
		offset += 17

		data[1] = 5
	}

	return BytesArray(c.exec(data[:offset]))
}

//Zrangebyscore
func (c *Conn) Zrangebyscore(key string, argvs ...interface{}) ([][]byte, error) {
	if argvs == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(argvs)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '3'
	data[offset+3] = '\r'
	data[offset+4] = '\n'
	data[offset+5] = 'Z'
	data[offset+6] = 'R'
	data[offset+7] = 'A'
	data[offset+8] = 'N'
	data[offset+9] = 'G'
	data[offset+10] = 'E'
	data[offset+11] = 'B'
	data[offset+12] = 'Y'
	data[offset+13] = 'S'
	data[offset+14] = 'C'
	data[offset+15] = 'O'
	data[offset+16] = 'R'
	data[offset+17] = 'E'
	data[offset+18] = '\r'
	data[offset+19] = '\n'
	offset += 20

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//argvs
	for i := 0; i < len(argvs); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(argvs[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Zrank
func (c *Conn) Zrank(key string, member interface{}) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '5'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'R'
	data[10] = 'A'
	data[11] = 'N'
	data[12] = 'K'
	data[13] = '\r'
	data[14] = '\n'
	offset := 15

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zrem
func (c *Conn) Zrem(key string, members ...interface{}) (int64, error) {
	data := c.buff.next(1024)

	size := len(members) + 2
	data[0] = '*'
	offset := 1
	offset += copy(data[offset:], c.intToBytes(size))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '4'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'Z'
	data[offset+5] = 'R'
	data[offset+6] = 'E'
	data[offset+7] = 'M'
	data[offset+8] = '\r'
	data[offset+9] = '\n'
	offset += 10

	//members
	for i := 0; i < len(members); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(members[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return Int64(c.exec(data[:offset]))
}

//Zremrangebyrank
func (c *Conn) Zremrangebyrank(key string, start, stop int64) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '5'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'Z'
	data[10] = 'R'
	data[11] = 'E'
	data[12] = 'M'
	data[13] = 'R'
	data[14] = 'A'
	data[15] = 'N'
	data[16] = 'G'
	data[17] = 'E'
	data[18] = 'B'
	data[19] = 'Y'
	data[20] = 'R'
	data[21] = 'A'
	data[22] = 'N'
	data[23] = 'K'
	data[24] = '\r'
	data[25] = '\n'
	offset := 26

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//stop
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(stop)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zremrangebyscore
func (c *Conn) Zremrangebyscore(key string, min, max int64) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '1'
	data[6] = '6'
	data[7] = '\r'
	data[8] = '\n'
	data[9] = 'Z'
	data[10] = 'R'
	data[11] = 'E'
	data[12] = 'M'
	data[13] = 'R'
	data[14] = 'A'
	data[15] = 'N'
	data[16] = 'G'
	data[17] = 'E'
	data[18] = 'B'
	data[19] = 'Y'
	data[20] = 'S'
	data[21] = 'C'
	data[22] = 'O'
	data[23] = 'R'
	data[24] = 'E'
	data[25] = '\r'
	data[26] = '\n'
	offset := 27

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//min
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(min)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//max
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(max)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zrevrange
func (c *Conn) Zrevrange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '4'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '9'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'R'
	data[10] = 'E'
	data[11] = 'V'
	data[12] = 'R'
	data[13] = 'A'
	data[14] = 'N'
	data[15] = 'G'
	data[16] = 'E'
	data[17] = '\r'
	data[18] = '\n'
	offset := 19

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//start
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(start)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//stop
	data[offset] = '$'
	offset++

	val = c.int64ToBytes(stop)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//withScores
	if withScores {
		data[offset] = '$'
		data[offset+1] = '1'
		data[offset+2] = '1'
		data[offset+3] = '\r'
		data[offset+4] = '\n'
		data[offset+5] = 'W'
		data[offset+6] = 'I'
		data[offset+7] = 'T'
		data[offset+8] = 'H'
		data[offset+9] = 'S'
		data[offset+10] = 'C'
		data[offset+11] = 'O'
		data[offset+12] = 'R'
		data[offset+13] = 'E'
		data[offset+14] = 'S'
		data[offset+15] = '\r'
		data[offset+16] = '\n'
		offset += 17

		data[1] = 5
	}

	return BytesArray(c.exec(data[:offset]))
}

//Zrevrangebyscore
func (c *Conn) Zrevrangebyscore(key string, argvs ...interface{}) ([][]byte, error) {
	if argvs == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(argvs)+2)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '1'
	data[offset+2] = '6'
	data[offset+3] = '\r'
	data[offset+4] = '\n'
	data[offset+5] = 'Z'
	data[offset+6] = 'R'
	data[offset+7] = 'E'
	data[offset+8] = 'V'
	data[offset+9] = 'R'
	data[offset+10] = 'A'
	data[offset+11] = 'N'
	data[offset+12] = 'G'
	data[offset+13] = 'E'
	data[offset+14] = 'B'
	data[offset+15] = 'Y'
	data[offset+16] = 'S'
	data[offset+17] = 'C'
	data[offset+18] = 'O'
	data[offset+19] = 'R'
	data[offset+20] = 'E'
	data[offset+21] = '\r'
	data[offset+22] = '\n'
	offset += 23

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//argvs
	for i := 0; i < len(argvs); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(argvs[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}

//Zrevrank
func (c *Conn) Zrevrank(key string, member interface{}) (int64, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '8'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'R'
	data[10] = 'E'
	data[11] = 'V'
	data[12] = 'R'
	data[13] = 'A'
	data[14] = 'N'
	data[15] = 'K'
	data[16] = '\r'
	data[17] = '\n'
	offset := 18

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return Int64(c.exec(data[:offset]))
}

//Zscore
func (c *Conn) Zscore(key string, member interface{}) (string, error) {
	data := c.buff.next(1024)

	data[0] = '*'
	data[1] = '3'
	data[2] = '\r'
	data[3] = '\n'

	//command
	data[4] = '$'
	data[5] = '6'
	data[6] = '\r'
	data[7] = '\n'
	data[8] = 'Z'
	data[9] = 'S'
	data[10] = 'C'
	data[11] = 'O'
	data[12] = 'R'
	data[13] = 'E'
	data[14] = '\r'
	data[15] = '\n'
	offset := 16

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//member
	data[offset] = '$'
	offset++

	val := c.interfaceToBytes(member)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	return String(c.exec(data[:offset]))
}

//Zscan
func (c *Conn) Zscan(key string, cursor int64, arguments ...interface{}) ([][]byte, error) {
	if arguments == nil {
		return nil, ErrNilKey
	}

	data := c.buff.next(1024)

	data[0] = '*'
	offset := copy(data[1:], c.intToBytes(len(arguments)+3)) + 1
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//command
	data[offset] = '$'
	data[offset+1] = '5'
	data[offset+2] = '\r'
	data[offset+3] = '\n'
	data[offset+4] = 'Z'
	data[offset+5] = 'S'
	data[offset+6] = 'C'
	data[offset+7] = 'A'
	data[offset+8] = 'N'
	data[offset+9] = '\r'
	data[offset+10] = '\n'
	offset += 11

	//key
	data[offset] = '$'
	offset++
	offset += copy(data[offset:], c.intToBytes(len(key)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], key)
	data[offset] = '\r'
	data[offset] = '\n'
	offset += 2

	//cursor
	data[offset] = '$'
	offset++

	val := c.int64ToBytes(cursor)
	offset += copy(data[offset:], c.intToBytes(len(val)))
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2
	offset += copy(data[offset:], val)
	data[offset] = '\r'
	data[offset+1] = '\n'
	offset += 2

	//values
	for i := 0; i < len(arguments); i++ {
		data[offset] = '$'
		offset++

		b := c.interfaceToBytes(arguments[i])
		offset += copy(data[offset:], c.intToBytes(len(b)))
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		offset += copy(data[offset:], b)
		data[offset] = '\r'
		data[offset+1] = '\n'
		offset += 2

		data = c.checkBuffer(offset, data)
	}

	return BytesArray(c.exec(data[:offset]))
}
