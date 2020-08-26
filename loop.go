package timer

import (
	"github.com/hslam/sortedlist"
	"sync"
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
	return &loopInstance{once: &sync.Once{}}
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
	once *sync.Once
	l    *loop
}

func getLoopInstance(instance *loopInstance, d time.Duration) *loop {
	instance.once.Do(func() {
		instance.l = newLoop(d)
		go func(instance *loopInstance) {
			idleTime := min(maxIdleTime, max(minIdleTime, instance.l.d*3))
			var lastIdle = time.Now()
			for {
				if len(instance.l.m) > 0 {
					lastIdle = time.Now()
				}
				if lastIdle.Add(idleTime).Before(time.Now()) {
					instance.l.Stop()
					instance.l = nil
					instance.once = &sync.Once{}
					break
				}
				Sleep(time.Second)
			}
		}(instance)
	})
	return instance.l
}

type loop struct {
	mu     sync.RWMutex
	sorted *sortedlist.SortedList
	m      map[*runtimeTimer]bool
	d      time.Duration
	done   chan struct{}
	closed bool
}

func newLoop(d time.Duration) *loop {
	l := &loop{
		d:      d,
		sorted: sortedlist.New(sortedlist.LessInt64),
		m:      make(map[*runtimeTimer]bool),
		done:   make(chan struct{}, 1),
	}
	go l.run()
	return l
}
func (l *loop) Register(r *runtimeTimer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.register(r)
}
func (l *loop) register(r *runtimeTimer) {
	l.m[r] = true
}
func (l *loop) Unregister(r *runtimeTimer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.unregister(r)
}
func (l *loop) unregister(r *runtimeTimer) {
	if _, ok := l.m[r]; ok {
		delete(l.m, r)
	}
}
func (l *loop) AddFunc(score int64, f timerFunc) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.addFunc(score, f)
}

func (l *loop) RunFunc(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.runFunc(now)
}

func (l *loop) addFunc(score int64, f timerFunc) {
	l.sorted.Insert(score, f)
}

func (l *loop) runFunc(now time.Time) bool {
	ss := make([]int64, 0)
	fs := make([]timerFunc, 0)
	for {
		if l.sorted.Length() < 1 {
			break
		}
		if l.sorted.Front().Score().(int64) > int64(now.UnixNano()) {
			break
		}
		top := l.sorted.Top()
		score, f := top.Value().(timerFunc)(now)
		if score > 0 {
			ss = append(ss, score)
			fs = append(fs, f)
		}
	}
	for i, score := range ss {
		l.addFunc(score, fs[i])
	}
	if len(ss) > 0 {
		return true
	} else {
		return false
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
	if l.closed {
		return true
	}
	l.closed = true
	close(l.done)
	l.m = nil
	return true
}
