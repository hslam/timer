// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"github.com/hslam/sortedlist"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	idleTime = time.Second
)

var (
	timersLen = int64(runtime.NumCPU())
	timers    = make([]timersBucket, timersLen)
)

type timer struct {
	closed int32
	when   int64
	period int64
	f      func(interface{})
	arg    interface{}
	event  func() (tick bool)
	bucket *timersBucket
}

func (r *timer) Start() {
	atomic.StoreInt32(&r.closed, 0)
	r.event = func() bool {
		if atomic.LoadInt32(&r.closed) == 0 {
			r.f(r.arg)
			if r.period > 0 {
				r.when += r.period
				return true
			}
		}
		return false
	}
	if r.bucket == nil {
		r.bucket = r.assignBucket()
	}
	r.bucket.GetInstance().AddTimer(r)
}

func (r *timer) assignBucket() *timersBucket {
	id := r.when / 1000 % timersLen
	return &timers[id]
}

func (r *timer) Stop() bool {
	if !atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		return true
	}
	if r.bucket != nil {
		r.bucket.GetInstance().DelTimer(r)
	}
	return true
}

type timersBucket struct {
	lock     sync.RWMutex
	created  bool
	sorted   *sortedlist.SortedList
	lastIdle time.Time
	pending  map[int64]struct{}
	trigger  chan struct{}
	wait     chan struct{}
	done     chan struct{}
	closed   int32
}

func (t *timersBucket) GetInstance() *timersBucket {
	t.lock.Lock()
	if !t.created {
		t.lastIdle = time.Now()
		t.created = true
		t.pending = make(map[int64]struct{})
		t.trigger = make(chan struct{}, 10)
		t.wait = make(chan struct{}, 10)
		t.sorted = sortedlist.New(sortedlist.LessInt64)
		t.done = make(chan struct{}, 1)
		atomic.StoreInt32(&t.closed, 0)
		t.lock.Unlock()
		go func(t *timersBucket) {
			for {
				Sleep(time.Second)
				t.lock.Lock()
				if t.lastIdle.Add(idleTime).Before(time.Now()) {
					if len(t.pending) == 0 && t.sorted.Length() == 0 {
						t.Stop()
						t.created = false
						t.lock.Unlock()
						break
					}
				}
				t.lock.Unlock()
			}
		}(t)
		go t.run()
	} else {
		t.lock.Unlock()
	}
	return t
}

func (t *timersBucket) AddTimer(r *timer) {
	t.lock.Lock()
	t.addTimer(r)
	when := t.sorted.Rear().Prev().Value().(*timer).when
	t.lastIdle = time.Unix(when/1000000000, when%1000000000)
	first := t.sorted.Front().Value().(*timer) == r
	t.lock.Unlock()
	if first {
		trigger(t.wait)
	}
	trigger(t.trigger)
}

func (t *timersBucket) DelTimer(r *timer) {
	t.lock.Lock()
	t.delTimer(r)
	t.lock.Unlock()
}

func (t *timersBucket) RunEvent(now time.Time) {
	t.lock.Lock()
	t.runEvent(now)
	t.lock.Unlock()
}

func (t *timersBucket) Front() int64 {
	t.lock.RLock()
	front := t.front()
	t.lock.RUnlock()
	return front
}

func (t *timersBucket) addTimer(r *timer) {
	t.sorted.Insert(r.when, r)
}

func (t *timersBucket) delTimer(r *timer) {
	t.sorted.Remove(r.when, r)
}

func (t *timersBucket) runEvent(now time.Time) {
	for {
		if t.sorted.Length() < 1 {
			break
		}
		if t.sorted.Front().Score().(int64) > int64(now.UnixNano()) {
			break
		}
		top := t.sorted.Top()
		r := top.Value().(*timer)
		if tick := r.event(); tick {
			t.addTimer(r)
		}
	}
}

func (t *timersBucket) front() int64 {
	if t.sorted.Length() < 1 {
		return 0
	}
	return t.sorted.Front().Score().(int64)
}

func (t *timersBucket) run() {
	for {
		when := t.Front()
		if when == 0 {
			select {
			case <-t.done:
				goto endfor
			case <-t.trigger:
				resetTrigger(t.trigger)
			}
			time.Sleep(time.Microsecond)
			continue
		}
		if when > runtimeNano() {
			t.lock.Lock()
			if _, ok := t.pending[when]; !ok {
				t.pending[when] = struct{}{}
				go func(when int64, t *timersBucket) {
					delta := time.Duration(when - runtimeNano())
					if delta > time.Microsecond {
						Sleep(delta)
					}
					select {
					case t.wait <- struct{}{}:
					default:
					}
					t.lock.Lock()
					delete(t.pending, when)
					t.lock.Unlock()
				}(when, t)
			}
			t.lock.Unlock()
			<-t.wait
			resetTrigger(t.wait)
		}
		t.RunEvent(time.Now())
	}
endfor:
}

func (t *timersBucket) Stop() bool {
	if atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		close(t.done)
	}
	return true
}

func trigger(c chan struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
}

func resetTrigger(c chan struct{}) {
	for len(c) > 0 {
		resetTriggerOnce(c)
	}
}

func resetTriggerOnce(c chan struct{}) {
	select {
	case <-c:
	default:
	}
}
