package redis

type Transaction struct {
	Conn

	status int
	count  int
}

const (
	WatchState = 1 << iota
	MultiState
	SubscribeState
	MonitorState
)

func (this *Conn) Transaction() *Transaction {
	return &Transaction{
		Conn: *(this),
	}
}

func (this *Transaction) Multi() error {
	var res, err = this.Do("MULTI")
	if err != nil {
		return err
	}

	if r, ok := res.([]byte); ok {
		if string(r) != "OK" {
			return ErrTransactionBegin
		} else {
			this.status |= MultiState
			return nil
		}
	}

	return ErrTransactionBegin
}

func (this *Transaction) Add(command string, args ...interface{}) error {
	var res, err = this.Do(command, args...)
	if err != nil {
		return err
	}

	if r, ok := res.([]byte); ok {
		if string(r) != "QUEUED" {
			return ErrTransactionAdd
		} else {
			this.count++
			return nil
		}
	}

	return ErrTransactionAdd
}

func (this *Transaction) Watch(keys ...interface{}) error {
	var res, err = this.Do("WATCH", keys...)
	if err != nil {
		return err
	}

	if r, ok := res.([]byte); ok {
		if string(r) != "OK" {
			return ErrWatchKey
		} else {
			this.status |= WatchState
			return nil
		}
	}

	return ErrWatchKey
}

func (this *Transaction) UnWatch() error {
	var res, err = this.Do("UNWATCH")
	if err != nil {
		return err
	}

	if r, ok := res.([]byte); ok {
		if string(r) != "OK" {
			return ErrUnWatchKey
		} else {
			return nil
		}
	}

	return ErrUnWatchKey
}

func (this *Transaction) Discard() error {
	var res, err = this.Do("DISCARD")
	if err != nil {
		return err
	}

	if r, ok := res.([]byte); ok {
		if string(r) != "OK" {
			return ErrDiscard
		} else {
			return nil
		}
	}

	return ErrDiscard
}

func (this *Transaction) Exec() (interface{}, error) {
	var res, err = this.Do("EXEC")
	if err != nil {
		return nil, err
	}

	return res, err
}
