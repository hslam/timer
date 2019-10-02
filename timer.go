package timer

import "time"



type runtimeTimer struct {
	arg    chan time.Time
	d 		time.Duration
	stop	chan bool
	closed	chan bool
	f 		func()
}

func (t *runtimeTimer) Start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		var startSleepTime time.Duration=0
		var lastSleepTime time.Duration=0
		var startWorkTime time.Duration=0
		var lastWorkTime time.Duration=0
		var d time.Duration=0
		var sd time.Duration=0
		for{
			startSleepTime=time.Duration(time.Now().UnixNano())
			select {
			case <-t.stop:
				goto endfor
			default:
				if lastSleepTime>t.d{
					d=t.d+time.Duration(float64(t.d-lastSleepTime)*BETA)
				}else {
					d=t.d
				}
				d-=lastWorkTime
				sd=time.Duration(ALPHA*float64(sd) + ((1 - ALPHA) * float64(d)))
				if sd>time.Microsecond{
					Sleep(sd)
				}else {
					Sleep(time.Microsecond)
				}
				lastSleepTime=time.Duration(time.Now().UnixNano())-startSleepTime

				startWorkTime=time.Duration(time.Now().UnixNano())
				if t.f!=nil{
					t.f()
				}else if t.arg!=nil{
					t.arg<-time.Now()
				}
				lastWorkTime=time.Duration(time.Now().UnixNano())-startWorkTime
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
		close(t.closed)
		return true
	case <-time.After(time.Second):
		close(t.closed)
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
	c := make(chan time.Time, 1)
	go func() {
		Sleep(d)
		c<-time.Now()
	}()
	return c
}