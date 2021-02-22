// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"errors"
	"sync/atomic"
	"time"
)

// A Ticker holds a channel that delivers `ticks' of a clock
// at intervals.
type Ticker struct {
	C      <-chan time.Time // The channel on which the ticks are delivered.
	r      timer
	closed int32
}

// NewTicker returns a new Ticker containing a channel that will send the
// time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will panic.
// Stop the ticker to release associated resources.
func NewTicker(d time.Duration) *Ticker {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	c := make(chan time.Time, 1)
	t := &Ticker{
		C: c,
		r: timer{
			when:   when(d),
			period: int64(d),
			f:      f,
			arg:    c,
		},
	}
	startTimer(&t.r)
	return t
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a concurrent goroutine
// reading from the channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return
	}
	stopTimer(&t.r)
	if c, ok := t.r.arg.(chan time.Time); ok {
		close(c)
	}
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. While Tick is useful for clients that have no need to shut down
// the Ticker, be aware that without a way to shut it down the underlying
// Ticker cannot be recovered by the garbage collector; it "leaks".
// Unlike NewTicker, Tick will return nil if d <= 0.
func Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
}

// TickFunc returns a new Ticker containing a func f that will call in its own goroutine
// with a period specified by the duration argument.
// It returns a Ticker that can be used to cancel the call using its Stop method.
func TickFunc(d time.Duration, f func()) *Ticker {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for TickFunc"))
	}
	type tickFunc struct {
		f     func()
		count int32
	}
	t := &Ticker{
		r: timer{
			when:   when(d),
			period: int64(d),
			f: func(arg interface{}) {
				tick := arg.(*tickFunc)
				if atomic.CompareAndSwapInt32(&tick.count, 0, 1) {
					go func(tick *tickFunc) {
						tick.f()
						atomic.StoreInt32(&tick.count, 0)
					}(tick)
				}
			},
			arg: &tickFunc{f: f},
		},
	}
	startTimer(&t.r)
	return t
}
