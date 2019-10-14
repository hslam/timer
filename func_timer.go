package timer

import(
	"time"
)

type funcTimer struct {
	d 		time.Duration
	when   	int64
	stop	chan bool
	closed	chan bool
	f 		func()
	count 	int64
}

func (f *funcTimer) Start() {
	go f.run()
}

func (f *funcTimer) run() {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	var s time.Duration=0
	var d time.Duration=0
	var sd time.Duration=0
	var lastSleepTime time.Duration=0
	var lastTime = time.Now()
	for{
		select {
		case <-f.stop:
			close(f.stop)
			goto endfor
		default:
			if lastSleepTime>f.d{
				d=lastSleepTime-f.d
			}else {
				d=f.d
			}
			sd = time.Duration(ALPHA*float64(sd) + ((1 - ALPHA) * float64(d)))
			s = min(f.d, max(time.Microsecond, time.Duration(BETA*float32(d))))
			Sleep(s)
			now:=time.Now()
			when:=f.when+f.count*int64(f.d)
			if when<int64(now.UnixNano()){
				f.count+=1
				if f.f!=nil{
					func() {
						defer func() {if err := recover(); err != nil {}}()
						f.f()
					}()
				}
			}
			lastSleepTime=time.Duration(now.UnixNano())-time.Duration(lastTime.UnixNano())
			lastTime=now
		}
	}
endfor:
	f.closed<-true
}

func (f *funcTimer) Stop() bool{
	defer func() {if err := recover(); err != nil {}}()
	defer func() {
		defer func() {if err := recover(); err != nil {}}()
		if f.closed!=nil{
			close(f.closed)
		}
	}()
	if f.stop==nil{
		return true
	}
	f.stop<-true
	select {
	case <-f.closed:
		return true
	case <-time.After(time.Second):
		return false
	}
}

func startFuncTimer(f *funcTimer){
	f.Start()
}

func stopFuncTimer(f *funcTimer) bool{
	return f.Stop()
}