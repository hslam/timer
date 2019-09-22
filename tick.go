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
			d:		d,
			arg:    c,
			stop:make(chan bool,1),
			closed:make(chan bool,1),
		},
	}
	startTimer(&t.r)
	return t
}
func (t *Ticker) TickFunc( f func()) {
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

