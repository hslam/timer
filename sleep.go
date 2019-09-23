package timer

/*
#include <unistd.h>
*/
import "C"
import (
	"time"
	"errors"
)
func Sleep(d time.Duration) {
	if d < time.Microsecond {
		panic(errors.New("non-positive interval for Sleep"))
	}
	var duration C.uint
	duration=C.uint(int64(d)/1000)
	C.usleep(duration)
}


type runtimeTimer struct {
	arg    chan time.Time
	d 		time.Duration
	stop	chan bool
	closed	chan bool
	f 		func()
}

func (t *runtimeTimer) Start() {
	go func() {
		var startTime time.Duration=0
		var lastWorkTime time.Duration=0
		var sd time.Duration=0
		for{
			select {
			case <-t.stop:
				goto endfor
			default:
				if lastWorkTime>0{
					sd=t.d-lastWorkTime
				}else {
					sd=t.d
				}
				if sd>time.Microsecond{
					Sleep(sd)
				}else {
					Sleep(time.Microsecond)
				}
				if t.f!=nil{
					startTime=time.Duration(time.Now().UnixNano())
					t.f()
					lastWorkTime=time.Duration(time.Now().UnixNano())-startTime
				}else if t.arg!=nil{
					t.arg<-time.Now()
				}
			}
		}
	endfor:
		t.closed<-true
	}()
}

func (t *runtimeTimer) Stop() bool{
	t.stop<-true
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