// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"errors"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	maxCount int64 = 0
	count    int64
)

func init() {
	if runtime.NumCPU() > 2 {
		maxCount = int64(runtime.NumCPU() - 2)
	}
}

func MaxCount(v int64) {
	maxCount = v
}

type Ticker struct {
	C <-chan time.Time
	r runtimeTimer
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

type FuncTicker struct {
	f        funcTimer
	r        runtimeTimer
	funcType bool
}

func NewFuncTicker(d time.Duration, f func()) *FuncTicker {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewFuncTicker"))
	}
	t := &FuncTicker{}
	if d <= time.Millisecond && atomic.LoadInt64(&count) < maxCount {
		atomic.AddInt64(&count, 1)
		t.funcType = true
		t.f = funcTimer{
			when: when(d),
			d:    d,
			f:    f,
			done: make(chan struct{}, 1),
		}
		startFuncTimer(&t.f)
		return t
	} else {
		t.r = runtimeTimer{
			tick:   true,
			when:   when(d),
			period: int64(d),
			f:      f,
		}
		startTimer(&t.r)
		return t
	}
}

func (t *FuncTicker) Tick(f func()) {
	if t.funcType {
		t.f.f = f
	} else {
		t.r.f = f
	}
}

func (t *FuncTicker) Stop() {
	if t.funcType {
		if !t.f.closed {
			atomic.AddInt64(&count, -1)
		}
		stopFuncTimer(&t.f)
	} else {
		stopTimer(&t.r)
	}
}

func (t *Ticker) Stop() {
	defer func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		if t.r.arg != nil {
			close(t.r.arg)
		}
	}()
	stopTimer(&t.r)
}

func Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
}
