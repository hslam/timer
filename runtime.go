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
		if atomic.LoadInt32(&r.closed) > 0 {
			return false
		}
		r.f(r.arg)
		if r.period > 0 && atomic.LoadInt32(&r.closed) == 0 {
			r.when += r.period
			return true
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
	lock     sync.Mutex
	created  bool
	mu       sync.RWMutex
	sorted   *sortedlist.SortedList
	lastIdle time.Time
	pmu      sync.Mutex
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
				if t.lastIdle.Add(idleTime).Before(time.Now()) && len(t.pending) == 0 {
					t.lock.Lock()
					t.Stop()
					t.created = false
					t.lock.Unlock()
					break
				}
				Sleep(time.Second)
			}
		}(t)
		go t.run()
	} else {
		t.lock.Unlock()
	}
	return t
}

func (t *timersBucket) Length() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.sorted.Length()
}

func (t *timersBucket) AddTimer(r *timer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.addTimer(r)
	when := t.sorted.Rear().Prev().Value().(*timer).when
	t.lastIdle = time.Unix(when/1000000000, when%1000000000)
	if t.sorted.Length() > 1 && t.sorted.Front().Value().(*timer) == r {
		select {
		case t.wait <- struct{}{}:
		default:
		}
	}
	select {
	case t.trigger <- struct{}{}:
	default:
	}
}

func (t *timersBucket) DelTimer(r *timer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.delTimer(r)
}

func (t *timersBucket) RunEvent(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.runEvent(now)
}

func (t *timersBucket) Front() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.front()
}

func (t *timersBucket) addTimer(r *timer) {
	t.sorted.Insert(r.when, r)
}

func (t *timersBucket) delTimer(r *timer) {
	t.sorted.Remove(r.when, r)
}

func (t *timersBucket) runEvent(now time.Time) {
	when := int64(0)
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
			when = r.when
		}
	}
	if when > 0 {
		t.lastIdle = time.Unix(when/1000000000, when%1000000000)
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
				if len(t.trigger) > 0 {
					<-t.trigger
				}
			}
			continue
		}
		if when > runtimeNano() {
			t.pmu.Lock()
			if _, ok := t.pending[when]; ok {
				t.pmu.Unlock()
			} else {
				t.pending[when] = struct{}{}
				t.pmu.Unlock()
				go func(when int64, t *timersBucket) {
					delta := time.Duration(when - runtimeNano())
					if delta > time.Microsecond {
						Sleep(delta)
					}
					select {
					case t.wait <- struct{}{}:
					default:
					}
					t.pmu.Lock()
					delete(t.pending, when)
					t.pmu.Unlock()
				}(when, t)
			}
			select {
			case <-t.done:
				goto endfor
			case <-t.wait:
				if len(t.wait) > 0 {
					<-t.wait
				}
			}
		}
		t.RunEvent(time.Now())
	}
endfor:
}

func (t *timersBucket) Stop() bool {
	if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return true
	}
	close(t.done)
	close(t.trigger)
	close(t.wait)
	return true
}
