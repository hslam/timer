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

type timerFunc func(now time.Time) (score int64, f timerFunc)

type runtimeTimer struct {
	arg       chan time.Time
	work      bool
	closed    int32
	when      int64
	period    int64
	f         func()
	timerFunc timerFunc
	count     int64
}

func (r *runtimeTimer) Start() {
	atomic.StoreInt32(&r.closed, 0)
	r.work = true
	r.count = 0
	r.timerFunc = func(now time.Time) (int64, timerFunc) {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		if atomic.LoadInt32(&r.closed) > 0 {
			return -1, nil
		}
		r.count += 1
		if r.f != nil && r.work {
			r.work = false
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
				}()
				r.f()
				r.work = true
			}()
		}
		if len(r.arg) == 0 {
			r.arg <- now
		}
		if r.period > 0 && atomic.LoadInt32(&r.closed) == 0 {
			return r.when + r.count*int64(r.period), r.timerFunc
		} else {
			return -1, nil
		}
	}
	getLoop(time.Duration(r.period)).AddFunc(r.when, r.timerFunc)
}

func (r *runtimeTimer) Stop() bool {
	if !atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		return true
	}
	return true
}

func startTimer(r *runtimeTimer) {
	r.Start()
}

func stopTimer(r *runtimeTimer) bool {
	return r.Stop()
}

func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

type Timer struct {
	C      <-chan time.Time
	r      runtimeTimer
	closed int32
}

func NewTimer(d time.Duration) *Timer {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTimer"))
	}
	c := make(chan time.Time, 1)
	r := &Timer{
		C: c,
		r: runtimeTimer{
			when: when(d),
			arg:  c,
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
	close(t.r.arg)
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
