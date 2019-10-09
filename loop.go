package timer

import (
	"time"
	"sync"
)

var (
	secondLoop 		= &loopInstance{}
	milliLoop100 	= &loopInstance{}
	milliLoop10 	= &loopInstance{}
	milliLoop 		= &loopInstance{}
	microLoop100 	= &loopInstance{}
	microLoop10 	= &loopInstance{}
	microLoop 		= &loopInstance{}
)

func GetLoop(d time.Duration) *loop {
	if d>=time.Second{
		return GetLoopInstance(secondLoop,time.Second)
	}else if d>=time.Millisecond*100{
		return GetLoopInstance(milliLoop100,time.Millisecond*100)
	}else if d>=time.Millisecond*10{
		return GetLoopInstance(milliLoop10,time.Millisecond*10)
	}else if d>=time.Millisecond{
		return GetLoopInstance(milliLoop,time.Millisecond)
	}else if d>=time.Microsecond*100{
		return GetLoopInstance(microLoop100,time.Microsecond*100)
	}else if d>=time.Microsecond*10{
		return GetLoopInstance(microLoop10,time.Microsecond*10)
	}else {
		return GetLoopInstance(microLoop,time.Microsecond)
	}
}

type loopInstance struct {
	once sync.Once
	l *loop
}

func GetLoopInstance(instance *loopInstance,d time.Duration)*loop  {
	instance.once.Do(func() {
		instance.l=newLoop(d)
	})
	return instance.l
}

type loop struct {
	mu 		sync.RWMutex
	sorted	*SortedList
	d 		time.Duration
	stop	chan bool
	closed	chan bool
}

func newLoop(d time.Duration) *loop {
	l:=&loop{
		d:			d,
		sorted:		NewSortedList(),
		stop:		make(chan bool,1),
		closed:		make(chan bool,1),
	}
	go l.run()
	return l
}

func (l *loop) AddFunc(score int64,f TimerFunc){
	l.mu.Lock()
	defer l.mu.Unlock()
	l.addFunc(score,f)
}

func (l *loop) RunFunc(now time.Time)bool{
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.runFunc(now)
}

func (l *loop) addFunc(score int64,f TimerFunc){
	l.sorted.Insert(Score(score),Value(f))
}

func (l *loop) runFunc(now time.Time)bool{
	ss:=make([]int64,0)
	fs:=make([]TimerFunc,0)
	for {
		if l.sorted.Length()<1{
			break
		}
		if int64(l.sorted.Rear().Score())>int64(now.UnixNano()){
			break
		}
		top:=l.sorted.Top()
		score,f:=top.Value()(now)
		if score>0{
			ss=append(ss,score)
			fs=append(fs,f)
		}
	}
	for i,score:=range ss{
		l.addFunc(score,fs[i])
	}
	if len(ss)>0{
		return true
	}else {
		return false
	}
}

func (l *loop) run() {
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
		case <-l.stop:
			goto endfor
		default:
			if lastSleepTime>l.d{
				d=lastSleepTime-l.d
			}else {
				d=l.d
			}
			sd = time.Duration(ALPHA*float64(sd) + ((1 - ALPHA) * float64(d)))
			s = min(l.d, max(time.Microsecond, time.Duration(BETA*float32(d))))
			Sleep(s)
			now:=time.Now()
			l.RunFunc(now)
			lastSleepTime=time.Duration(now.UnixNano())-time.Duration(lastTime.UnixNano())
			lastTime=now
		}
	}
endfor:
	l.closed<-true
}

func (l *loop) Stop() bool{
	defer func() {
		l=nil
	}()
	l.stop<-true
	select {
	case <-l.closed:
		close(l.closed)
		return true
	case <-time.After(time.Second):
		return false
	}
}
