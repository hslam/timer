package timer

import (
	"time"
	"errors"
	"sync/atomic"
	"runtime"
)

var(
	maxCount 	int64 =0
	count 		int64
)

func init() {
	if runtime.NumCPU()>2{
		maxCount=int64(runtime.NumCPU()-2)
	}
}
type Ticker struct {
	C <-chan time.Time
	r runtimeTimer
}

func NewTicker(d time.Duration) *Ticker{
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	c := make(chan time.Time, 1)
	t := &Ticker{
		C: c,
		r: runtimeTimer{
			tick:true,
			when: when(d),
			period: int64(d),
			arg:    c,
			closed:make(chan bool,1),
		},
	}
	startTimer(&t.r)
	return t
}

type FuncTicker struct {
	Ticker
	f funcTimer
	funcType bool
}

func NewFuncTicker(d time.Duration,f func()) *FuncTicker{
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	t := &FuncTicker{}
	if d<=time.Millisecond&&atomic.LoadInt64(&count)<maxCount{
		atomic.AddInt64(&count,1)
		t.funcType=true
		t.f= funcTimer{
			when: when(d),
			d:		d,
			f:		f,
			stop:make(chan bool,1),
			closed:make(chan bool,1),
		}
		startFuncTimer(&t.f)
		return t
	}else {
		t.r= runtimeTimer{
			tick:true,
			when: when(d),
			period: int64(d),
			f:		f,
			closed:make(chan bool,1),
		}
		startTimer(&t.r)
		return t
	}
}

func (t *FuncTicker) Tick( f func()) {
	 if t.funcType{
		t.f.f=f
	}else{
		t.r.f=f
	}
}

func (t *FuncTicker) Stop() {
	if t.funcType{
		atomic.AddInt64(&count,-1)
		stopFuncTimer(&t.f)
	}else {
		stopTimer(&t.r)
	}
}

func (t *Ticker) Stop() {
	defer func() {
		defer func() {if err := recover(); err != nil {}}()
		if t.r.arg!=nil{
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

