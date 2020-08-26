package timer

import (
	"math"
	"testing"
	"time"
)

func TestTimerSleep(t *testing.T) {
	startTime := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond)
	}
	runTime := time.Now().UnixNano() - startTime
	startTimer := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		Sleep(time.Millisecond)
	}
	runTimer := time.Now().UnixNano() - startTimer
	ratioTimer := float64(runTimer-int64(time.Second)) / float64(time.Second)
	ratioTime := float64(runTime-int64(time.Second)) / float64(time.Second)
	if ratioTimer > ratioTime && ratioTimer-ratioTime > 0.1 {
		t.Errorf("error ratio Timer:%.2f Time:%.2f", ratioTimer, ratioTime)
	}
	t.Log("timer.Sleep()", runTimer, ",time.Sleep()", runTime, ",ratioTimer-ratioTime", ratioTimer-ratioTime)
}

func BenchmarkTimerSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sleep(time.Millisecond)
	}
}

func TestTimerTicker(t *testing.T) {
	timeTicker := time.NewTicker(time.Millisecond)
	startTime := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		<-timeTicker.C
	}
	runTime := time.Now().UnixNano() - startTime
	defer timeTicker.Stop()

	timerTicker := NewTicker(time.Millisecond)
	startTimer := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		<-timerTicker.C
	}
	runTimer := time.Now().UnixNano() - startTimer
	defer timerTicker.Stop()
	ratioTimer := float64(runTimer-int64(time.Second)) / float64(time.Second)
	ratioTime := float64(runTime-int64(time.Second)) / float64(time.Second)
	if math.Abs(ratioTimer-ratioTime) > 0.1 || math.Abs(ratioTimer) > 0.1 {
		t.Errorf("error ratio Timer:%.2f Time:%.2f", ratioTimer, ratioTime)
	}
	t.Log("timer.Ticker", runTimer, ",time.Ticker", runTime, ",ratioTimer", ratioTimer, ",ratioTime", ratioTime)
}

//go test -v -bench=. -benchmem
//=== RUN   TestTimerSleep
//--- PASS: TestTimerSleep (2.88s)
//timer_test.go:25: timer.Sleep() 1461301000 ,time.Sleep() 1413659000 ,ratioTimer-ratioTime 0.04764200000000002
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.00s)
//timer_test.go:55: timer.Ticker 1001063000 ,time.Ticker 1000397000 ,ratioTimer 0.001063 ,ratioTime 0.000397
//goos: darwin
//goarch: amd64
//pkg: github.com/hslam/timer
//BenchmarkTimerSleep-4   	    1000	   1359823 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	github.com/hslam/timer	6.373s

//go test -v -bench=. -benchmem  -tags=use_cgo
//=== RUN   TestTimerSleep
//--- PASS: TestTimerSleep (2.79s)
//timer_test.go:25: timer.Sleep() 1358244000 ,time.Sleep() 1430921000 ,ratioTimer-ratioTime -0.07267699999999999
//=== RUN   TestTimerTicker
//--- PASS: TestTimerTicker (2.00s)
//timer_test.go:55: timer.Ticker 1001974000 ,time.Ticker 1000416000 ,ratioTimer 0.001974 ,ratioTime 0.000416
//goos: darwin
//goarch: amd64
//pkg: github.com/hslam/timer
//BenchmarkTimerSleep-4   	    1000	   1387910 ns/op	       0 B/op	       0 allocs/op
//PASS
//ok  	github.com/hslam/timer	6.329s
