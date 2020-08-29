// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package timer provides functionality for measuring time.
package timer

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Second                    = 1000 * Millisecond
	Minute                    = 60 * Second
	Hour                      = 60 * Minute
)

func startTimer(r *timer) {
	r.Start()
}

func stopTimer(r *timer) bool {
	return r.Stop()
}

func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

func AfterFunc(d time.Duration, f func()) *Timer {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for AfterFunc"))
	}
	t := &Timer{
		r: timer{
			when: when(d),
			f: func(arg interface{}) {
				go arg.(func())()
			},
			arg: f,
		},
	}
	startTimer(&t.r)
	return t
}

type Timer struct {
	C      <-chan time.Time
	r      timer
	closed int32
}

func NewTimer(d time.Duration) *Timer {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTimer"))
	}
	c := make(chan time.Time, 1)
	r := &Timer{
		C: c,
		r: timer{
			when: when(d),
			f: func(arg interface{}) {
				select {
				case arg.(chan time.Time) <- time.Now():
				default:
				}
			},
			arg: c,
		},
	}
	startTimer(&r.r)
	return r
}

func (t *Timer) Stop() bool {
	if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return true
	}
	active := stopTimer(&t.r)
	if c, ok := t.r.arg.(chan time.Time); ok {
		close(c)
	}
	return active
}

func (t *Timer) Reset(d time.Duration) bool {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for Reset"))
	}
	w := when(d)
	active := stopTimer(&t.r)
	t.r.when = w
	startTimer(&t.r)
	return active
}

func when(d time.Duration) int64 {
	if d <= 0 {
		return runtimeNano()
	}
	r := runtimeNano() + int64(d)
	if r < 0 {
		r = 1<<63 - 1 // math.MaxInt64
	}
	return r
}

func runtimeNano() int64 {
	return time.Now().UnixNano()
}
