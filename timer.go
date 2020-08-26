// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"errors"
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

type timerFunc func(now time.Time) (score int64, r *runtimeTimer)

type runtimeTimer struct {
	tick      bool
	arg       chan time.Time
	work      bool
	closed    bool
	when      int64
	period    int64
	f         func()
	timerFunc timerFunc
	count     int64
}

func (r *runtimeTimer) Start() {
	r.work = true
	r.timerFunc = func(now time.Time) (int64, *runtimeTimer) {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		if r.closed {
			return -1, r
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
		} else if r.arg != nil && len(r.arg) == 0 {
			func() {
				defer func() {
					if err := recover(); err != nil {
					}
				}()
				r.arg <- now
			}()
		}
		if r.tick && !r.closed {
			return r.when + r.count*int64(r.period), r
		} else {
			return -1, r
		}
	}
	getLoop(time.Duration(r.period)).Register(r)
	getLoop(time.Duration(r.period)).AddFunc(r.when, r.timerFunc)
}

func (r *runtimeTimer) Stop() bool {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if r.closed {
		return true
	}
	r.closed = true
	getLoop(time.Duration(r.period)).Unregister(r)
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
	C <-chan time.Time
	r runtimeTimer
}

func NewTimer(d time.Duration) *Timer {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTimer"))
	}
	c := make(chan time.Time, 1)
	r := &Timer{
		C: c,
		r: runtimeTimer{
			when:   when(d),
			period: int64(d),
			arg:    c,
		},
	}
	startTimer(&r.r)
	return r
}

func (t *Timer) Stop() bool {
	defer func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		if t.r.arg != nil {
			close(t.r.arg)
		}
	}()
	return stopTimer(&t.r)
}

func (t *Timer) Reset(d time.Duration) bool {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for Reset"))
	}
	w := when(d)
	r := t.r
	active := stopTimer(&r)
	t.r = runtimeTimer{
		when:   w,
		period: int64(d),
		arg:    r.arg,
	}
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
