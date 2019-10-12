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

func (r *runtimeTimer) Start() {
	r.work=true
	r.stop=false
	if r.closed==nil{
		r.closed=make(chan bool,1)
	}
	if r.workchan!=nil{
		go func() {
			for range r.workchan{
				if r.f!=nil{
					func(){
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						r.f()
					}()
				}
			}
		}()
	}
	r.timerFunc= func(now time.Time)(score int64,f timerFunc) {
		defer func() {if err := recover(); err != nil {}}()
		r.count+=1
		if r.f!=nil&&r.work{
			r.work=false
			r.workchan<-true
		}else if r.arg!=nil&&len(r.arg)==0{
			func() {
				defer func() {if err := recover(); err != nil {}}()
				r.arg<-now
			}()
		}
		if r.tick&&!r.stop{
			return r.when+r.count*int64(r.period),r.timerFunc
		}else {
			func() {
				defer func() {if err := recover(); err != nil {}}()
				r.closed<-true
			}()
			return -1,nil
		}
	}
	getLoop(time.Duration(r.period)).AddFunc(r.when,r.timerFunc)
}

func (r *runtimeTimer) Stop() bool{
	defer func() {if err := recover(); err != nil {}}()
	defer func() {
		defer func() {if err := recover(); err != nil {}}()
		if r.workchan!=nil{
			close(r.workchan)
		}
	}()
	defer func() {
		defer func() {if err := recover(); err != nil {}}()
		if r.closed!=nil{
			close(r.closed)
		}
	}()
	r.stop=true
	select {
	case <-r.closed:
		return true
	case <-time.After(time.Second):
		return false
	}
}

func startTimer(r *runtimeTimer){
	r.Start()
}

func stopTimer(r *runtimeTimer) bool{
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
			when: when(d),
			period: int64(d),
			arg:  c,
			closed:make(chan bool,1),
		},
	}
	startTimer(&r.r)
	return r
}

func (r *Timer) Stop() bool {
	defer func() {
		defer func() {if err := recover(); err != nil {}}()
		if r.r.arg!=nil{
			close(r.r.arg)
		}
	}()
	return stopTimer(&r.r)
}

func (r *Timer) Reset(d time.Duration) bool {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	w := when(d)
	active := stopTimer(&r.r)
	r.r.when = w
	r.r.period=int64(d)
	startTimer(&r.r)
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

func runtimeNano() int64{
	return time.Now().UnixNano()
}
