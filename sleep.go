package timer

/*
#include <unistd.h>
*/
import "C"
import (
	"time"
)
func Sleep(d time.Duration) {
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
		for{
			select {
			case <-t.stop:
				goto endfor
			default:
				Sleep(t.d)
				if t.f!=nil{
					t.f()
				}else {
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