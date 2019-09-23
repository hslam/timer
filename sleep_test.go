package timer

import (
	"testing"
	"time"
	"math"
)

func TestTimerSleep(t *testing.T) {
	startTime:=time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond)
	}
	runTime:=time.Now().UnixNano()-startTime

	startTimer:=time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		Sleep(time.Millisecond)
	}
	runTimer:=time.Now().UnixNano()-startTimer
	ratioTimer:=float64(runTimer-int64(time.Second))/float64(time.Second)
	ratioTime:=float64(runTime-int64(time.Second))/float64(time.Second)
	if ratioTimer>ratioTime&&ratioTimer-ratioTime>0.1{
		t.Errorf("error ratio Timer:%.2f Time:%.2f",ratioTimer,ratioTime)
	}
	t.Log("timer.Sleep()",runTimer,",time.Sleep()",runTime,",ratioTimer-ratioTime",ratioTimer-ratioTime)
}


func TestTimerTicker(t *testing.T) {
	timeTicker:=time.NewTicker(time.Millisecond)
	startTime:=time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		<-timeTicker.C
	}
	runTime:=time.Now().UnixNano()-startTime
	timeTicker.Stop()

	timerTicker:=NewTicker(time.Millisecond)
	startTimer:=time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		<-timerTicker.C
	}
	runTimer:=time.Now().UnixNano()-startTimer
	timerTicker.Stop()
	ratioTimer:=float64(runTimer-int64(time.Second))/float64(time.Second)
	ratioTime:=float64(runTime-int64(time.Second))/float64(time.Second)
	if math.Abs(ratioTimer-ratioTime)>0.1||math.Abs(ratioTimer)>0.1{
		t.Errorf("error ratio Timer:%.2f Time:%.2f",ratioTimer,ratioTime)
	}
	t.Log("timer.Ticker",runTimer,",time.Ticker",runTime,",ratioTimer",ratioTimer,",ratioTime",ratioTime)
}



func BenchmarkTimerSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sleep(time.Millisecond)
	}
}

func BenchmarkTimeSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond)
	}
}

/*
go test -v -bench=. -benchmem
=== RUN   TestTimerSleep
--- PASS: TestTimerSleep (2.81s)
    sleep_test.go:26: timer.Sleep() 1361899000 ,time.Sleep() 1445361000 ,ratioTimer-ratioTime -0.08346199999999998
=== RUN   TestTimerTicker
--- PASS: TestTimerTicker (2.03s)
    sleep_test.go:51: timer.Ticker 1026243000 ,time.Ticker 1000392000 ,ratioTimer 0.026243 ,ratioTime 0.000392
goos: darwin
goarch: amd64
pkg: hslam.com/mgit/Mort/timer
BenchmarkTimerSleep-4   	    1000	   1415950 ns/op	       0 B/op	       0 allocs/op
BenchmarkTimeSleep-4    	    1000	   1357983 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	hslam.com/mgit/Mort/timer	7.909s*/
