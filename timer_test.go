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
//--- PASS: TestTimerSleep (2.90s)
//timer_test.go:27: timer.Sleep() 1466564000 ,time.Sleep() 1434900000 ,ratioTimer-ratioTime 0.03166399999999997
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.00s)
//timer_test.go:57: timer.Ticker 1002659000 ,time.Ticker 1000402000 ,ratioTimer 0.002659 ,ratioTime 0.000402
//goos: darwin
//goarch: amd64
//pkg: hslam.com/mgit/Mort/timer
//BenchmarkTimerSleep-4   	    1000	   1465932 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	hslam.com/mgit/Mort/timer	6.515s

//go test -v -bench=. -benchmem  -tags=use_cgo
//=== RUN   TestTimerSleep
//--- PASS: TestTimerSleep (2.86s)
//timer_test.go:27: timer.Sleep() 1401094000 ,time.Sleep() 1459930000 ,ratioTimer-ratioTime -0.058836
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.01s)
//timer_test.go:57: timer.Ticker 1006298000 ,time.Ticker 1000383000 ,ratioTimer 0.006298 ,ratioTime 0.000383
//goos: darwin
//goarch: amd64
//pkg: hslam.com/mgit/Mort/timer
//BenchmarkTimerSleep-4   	    1000	   1409320 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	hslam.com/mgit/Mort/timer	6.430s

