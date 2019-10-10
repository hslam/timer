package timer

import (
	"time"
	"errors"
)

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
}

func NewFuncTicker(d time.Duration,f func()) *FuncTicker{
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	t := &FuncTicker{}
	t.r= runtimeTimer{
		tick:true,
		when: when(d),
		period: int64(d),
		f:		f,
		closed:make(chan bool,1),
		workchan:make(chan bool,1),
	}
	go func() {
		for range t.r.workchan{
			if t.r.f!=nil{
				func(){
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					t.r.f()
				}()
			}
			t.r.work=true
		}
	}()
	startTimer(&t.r)
	return t
}

func (t *FuncTicker) Tick( f func()) {
	t.r.f=f
}

func (t *Ticker) Stop() {
	stopTimer(&t.r)
}

func Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
}

