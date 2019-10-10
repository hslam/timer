package timer

import (
	"time"
	"errors"
)

const (
	Nanosecond  		time.Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

type timerFunc  func(now time.Time)(score int64,f timerFunc)

type runtimeTimer struct {
	tick 	bool
	arg		chan time.Time
	work 	bool
	stop	bool
	closed	chan bool
	when   	int64
	period 	int64
	f 		func()
	timerFunc	timerFunc
	workchan chan bool
	count int64
}

func (t *runtimeTimer) Start() {
	t.work=true
	t.stop=false
	if t.closed==nil{
		t.closed=make(chan bool,1)
	}
	t.timerFunc= func(now time.Time)(score int64,f timerFunc) {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		t.count+=1
		if t.f!=nil&&t.work{
			t.work=false
			t.workchan<-true
		}else if t.arg!=nil&&len(t.arg)==0{
			t.arg<-now
		}
		if t.tick&&!t.stop{
			return t.when+t.count*int64(t.period),t.timerFunc
		}else {
			t.closed<-true
			return -1,nil
		}
	}
	GetLoop(time.Duration(t.period)).AddFunc(t.when,t.timerFunc)
}

func (t *runtimeTimer) Stop() bool{
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	defer func() {
		close(t.closed)
	}()
	t.stop=true
	select {
	case <-t.closed:
		return true
	case <-time.After(time.Second):
		return false
	}
}

func startTimer(t *runtimeTimer){
	t.Start()
}

func stopTimer(t *runtimeTimer) bool{
	return t.Stop()
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
	t := &Timer{
		C: c,
		r: runtimeTimer{
			when: when(d),
			period: int64(d),
			arg:  c,
			closed:make(chan bool,1),
		},
	}
	startTimer(&t.r)
	return t
}

func (t *Timer) Stop() bool {
	return stopTimer(&t.r)
}

func (t *Timer) Reset(d time.Duration) bool {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	w := when(d)
	active := stopTimer(&t.r)
	t.r.when = w
	t.r.period=int64(d)
	startTimer(&t.r)
	return active
}

func when(d time.Duration) int64 {
	if d <= 0 {
		return runtimeNano()
	}
	t := runtimeNano() + int64(d)
	if t < 0 {
		t = 1<<63 - 1 // math.MaxInt64
	}
	return t
}

func runtimeNano() int64{
	return time.Now().UnixNano()
}
