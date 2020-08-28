// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"github.com/hslam/sortedlist"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxIdleTime = time.Second * 3
	minIdleTime = time.Millisecond * 30
)

var (
	secondLoop   = newloopInstance()
	milliLoop100 = newloopInstance()
	milliLoop10  = newloopInstance()
	milliLoop    = newloopInstance()
	microLoop100 = newloopInstance()
	microLoop10  = newloopInstance()
	microLoop    = newloopInstance()
)

func newloopInstance() *loopInstance {
	return &loopInstance{}
}

func getLoop(d time.Duration) *loop {
	for {
		l := selectLoop(d)
		if l != nil {
			return l
		}
	}
}

func selectLoop(d time.Duration) *loop {
	if d >= time.Second || d == 0 {
		return getLoopInstance(secondLoop, time.Second)
	} else if d >= time.Millisecond*100 {
		return getLoopInstance(milliLoop100, time.Millisecond*100)
	} else if d >= time.Millisecond*10 {
		return getLoopInstance(milliLoop10, time.Millisecond*10)
	} else if d >= time.Millisecond {
		return getLoopInstance(milliLoop, time.Millisecond)
	} else if d >= time.Microsecond*100 {
		return getLoopInstance(microLoop100, time.Microsecond*100)
	} else if d >= time.Microsecond*10 {
		return getLoopInstance(microLoop10, time.Microsecond*10)
	} else {
		return getLoopInstance(microLoop, time.Microsecond)
	}
}

type loopInstance struct {
	created int32
	l       *loop
}

func getLoopInstance(instance *loopInstance, d time.Duration) *loop {
	if atomic.CompareAndSwapInt32(&instance.created, 0, 1) {
		l := newLoop(d)
		instance.l = l
		go func(instance *loopInstance, l *loop) {
			var idleTime = min(maxIdleTime, max(minIdleTime, l.d*3))
			var lastIdle = time.Now()
			for {
				if l.Length() > 0 {
					lastIdle = time.Now()
				} else if lastIdle.Add(idleTime).Before(time.Now()) {
					atomic.StoreInt32(&instance.created, 0)
					l.Stop()
					break
				}
				Sleep(time.Second)
			}
		}(instance, l)
	}
	return instance.l
}

type loop struct {
	mu     sync.RWMutex
	sorted *sortedlist.SortedList
	d      time.Duration
	done   chan struct{}
	closed int32
}

func newLoop(d time.Duration) *loop {
	l := &loop{
		d:      d,
		sorted: sortedlist.New(sortedlist.LessInt64),
		done:   make(chan struct{}, 1),
	}
	go l.run()
	return l
}

func (l *loop) Length() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.sorted.Length()
}

func (l *loop) AddTimer(r *runtimeTimer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.addTimer(r)
}

func (l *loop) RunFunc(now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.runFunc(now)
}

func (l *loop) addTimer(r *runtimeTimer) {
	l.sorted.Insert(r.when, r)
}

func (l *loop) runFunc(now time.Time) {
	for {
		if l.sorted.Length() < 1 {
			break
		}
		if l.sorted.Front().Score().(int64) > int64(now.UnixNano()) {
			break
		}
		top := l.sorted.Top()
		r := top.Value().(*runtimeTimer)
		if r.timerFunc(now) {
			l.addTimer(r)
		}
	}
}

func (l *loop) run() {
	var s time.Duration = 0
	var d time.Duration = 0
	var sd time.Duration = 0
	var lastSleepTime time.Duration = 0
	var lastTime = time.Now()
	for {
		select {
		case <-l.done:
			goto endfor
		default:
			if lastSleepTime > l.d {
				d = lastSleepTime - l.d
			} else {
				d = l.d
			}
			sd = time.Duration(ALPHA*float64(sd) + ((1 - ALPHA) * float64(d)))
			s = min(l.d, max(time.Microsecond, time.Duration(BETA*float32(d))))
			Sleep(s)
			now := time.Now()
			l.RunFunc(now)
			lastSleepTime = time.Duration(now.UnixNano()) - time.Duration(lastTime.UnixNano())
			lastTime = now
		}
	}
endfor:
}

func (l *loop) Stop() bool {
	if !atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		return true
	}
	close(l.done)
	return true
}

func min(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

func max(a, b time.Duration) time.Duration {
	if a >= b {
		return a
	}
	return b
}
