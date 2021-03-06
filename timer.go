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
	// Nanosecond represents the nano second time.Duration.
	Nanosecond time.Duration = 1
	// Microsecond represents the micro second time.Duration.
	Microsecond = 1000 * Nanosecond
	// Millisecond represents the milli second time.Duration.
	Millisecond = 1000 * Microsecond
	// Second represents the second time.Duration.
	Second = 1000 * Millisecond
	// Minute represents the minute time.Duration.
	Minute = 60 * Second
	// Hour represents the hour time.Duration.
	Hour = 60 * Minute
)

func startTimer(r *timer) {
	r.Start()
}

func stopTimer(r *timer) bool {
	return r.Stop()
}

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to NewTimer(d).C.
// The underlying Timer is not recovered by the garbage collector
// until the timer fires. If efficiency is a concern, use NewTimer
// instead and call Timer.Stop if the timer is no longer needed.
func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

// AfterFunc waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer that can
// be used to cancel the call using its Stop method.
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

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	C      <-chan time.Time
	r      timer
	closed int32
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d time.Duration) *Timer {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTimer"))
	}
	c := make(chan time.Time, 1)
	r := &Timer{
		C: c,
		r: timer{
			when: when(d),
			f:    f,
			arg:  c,
		},
	}
	startTimer(&r.r)
	return r
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
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

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
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

func f(arg interface{}) {
	select {
	case arg.(chan time.Time) <- time.Now():
	default:
	}
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
