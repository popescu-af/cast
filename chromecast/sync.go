package chromecast

import (
	"fmt"
	"sync"
	"time"
)

type commandSync struct {
	ch      chan struct{}
	mtx     sync.RWMutex
	waiting bool
}

func newCommandSync() *commandSync {
	return &commandSync{
		ch: make(chan struct{}),
	}
}

func (c *commandSync) setWaiting(waiting bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.waiting = waiting
}

func (c *commandSync) wait(d time.Duration) error {
	c.mtx.RLock()
	waiting := c.waiting
	c.mtx.RUnlock()

	if !waiting {
		return fmt.Errorf("not waiting for command sync")
	}

	select {
	case <-c.ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timed out waiting for command sync")
	}
}

func (c *commandSync) notifyDone() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.waiting {
		c.ch <- struct{}{}
		c.waiting = false
	}
}
