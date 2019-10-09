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

func BenchmarkTimerSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sleep(time.Millisecond)
	}
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

//go test -v -bench=. -benchmem
//=== RUN   TestTimerSleep
//--- PASS: TestTimerSleep (2.66s)
//timer_test.go:27: timer.Sleep() 1335369000 ,time.Sleep() 1319799000 ,ratioTimer-ratioTime 0.015569999999999973
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.01s)
//timer_test.go:57: timer.Ticker 1005705000 ,time.Ticker 1000190000 ,ratioTimer 0.005705 ,ratioTime 0.00019
//goos: darwin
//goarch: amd64
//pkg: hslam.com/mgit/Mort/timer
//BenchmarkTimerSleep-4   	    1000	   1330641 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	hslam.com/mgit/Mort/timer	6.138s

//go test -v -bench=. -benchmem  -tags=use_cgo
//=== RUN   TestTimerSleep
//--- PASS: TestTimerSleep (2.67s)
//timer_test.go:27: timer.Sleep() 1299282000 ,time.Sleep() 1368566000 ,ratioTimer-ratioTime -0.06928400000000001
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.00s)
//timer_test.go:57: timer.Ticker 1001452000 ,time.Ticker 1000327000 ,ratioTimer 0.001452 ,ratioTime 0.000327
//goos: darwin
//goarch: amd64
//pkg: hslam.com/mgit/Mort/timer
//BenchmarkTimerSleep-4   	    1000	   1314661 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	hslam.com/mgit/Mort/timer	6.128s

