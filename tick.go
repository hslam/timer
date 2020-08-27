// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"errors"
	"sync/atomic"
	"time"
)

type Ticker struct {
	C      <-chan time.Time
	r      runtimeTimer
	closed int32
}

func NewTicker(d time.Duration) *Ticker {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	c := make(chan time.Time, 1)
	t := &Ticker{
		C: c,
		r: runtimeTimer{
			tick:   true,
			when:   when(d),
			period: int64(d),
			arg:    c,
		},
	}
	startTimer(&t.r)
	return t
}

func (t *Ticker) Tick(f func()) {
	t.r.f = f
}

func (t *Ticker) Stop() {
	if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return
	}
	stopTimer(&t.r)
	close(t.r.arg)
}

func Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
}
