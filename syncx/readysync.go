package helper

import (
	"fmt"
	"sync"
	"time"
)

/*Ready sync*/
type ReadySync struct {
	once sync.Once
	ch   chan bool
}

func ReadyContext() *ReadySync {
	return &ReadySync{
		ch: make(chan bool),
	}
}

func (r *ReadySync) Close() {
	close(r.ch)
}

func (r *ReadySync) Ready() {
	r.ch <- true
}

func (r *ReadySync) ReadyOnce() {
	r.once.Do(func() {
		r.Ready()
	})
}

func (r *ReadySync) Wait(t time.Duration) error {
	select {
	case <-time.After(t):
		return fmt.Errorf("Timeout after %v", t)
	case ok, more := <-r.ch:
		if !more || !ok {
			return fmt.Errorf("Not ready")
		}
	}
	return nil
}
