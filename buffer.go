package redis

import (
	"errors"
	"sync/atomic"
)

type buffer struct {
	freeCache chan *item
	min       int32
	max       int32
	active    int32
}

type item struct {
	buff []byte
	idx  int
	free int
}

//1.5 kb
var defaultBlockSize = 1536

func newBuffer(min, max int) *buffer {
	b := &buffer{
		freeCache: make(chan *item, max),
		min:       int32(min),
		max:       int32(max),
	}

	for i := int32(0); i < b.min; i++ {
		b.freeCache <- newItem()
		b.active++
	}

	return b
}

func newItem() *item {
	return &item{
		buff: make([]byte, defaultBlockSize),
		free: defaultBlockSize,
	}
}

func (b *buffer) get() *item {
	select {
	case t := <-b.freeCache:
		atomic.AddInt32(&b.active, -1)
		return t

	default:
		return newItem()
	}
}

func (b *buffer) put(t *item) {
	if b.active < b.max {
		b.freeCache <- t
		atomic.AddInt32(&b.active, 1)
	}
}

func (t *item) next(need int) []byte {
	if need > t.free {
		t.grow(need)
	}

	offset := t.idx
	t.idx += need
	t.free -= need
	return t.buff[offset:t.idx]
}

var errTooLarge = errors.New("cache to back is too large")

//attention: array will rewrite
func (t *item) back(need int) error {
	if need > t.idx {
		return errTooLarge
	}

	t.idx -= need
	t.free += need
	return nil
}

func (t *item) reset() {
	t.idx = 0
	t.free = len(t.buff)
}

func (t *item) grow(need int) {
	multiple := need/len(t.buff) + 1
	if multiple < 2 {
		multiple = 2
	}

	newBuf := make([]byte, multiple*len(t.buff))
	copy(newBuf, t.buff)
	t.buff = newBuf
	t.free = len(newBuf) - t.idx
}

func (t *item) data() []byte {
	return t.buff[:t.idx]
}
