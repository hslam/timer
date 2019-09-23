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
	ratio:=float64(runTimer-runTime)/float64(runTime)
	if ratio>0.1{
		t.Errorf("error ratio %.2f%%",ratio)
	}
	t.Log("timer.Sleep()",runTimer,",time.Sleep()",runTime,"ratio",ratio)
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
	ratio:=float64(runTimer-runTime)/float64(runTime)
	if math.Abs(ratio)>0.1{
		t.Errorf("error ratio %.2f%%",ratio)
	}
	t.Log("timer.Ticker",runTimer,",time.Ticker",runTime,"ratio",ratio)
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
goos: darwin
goarch: amd64
pkg: hslam.com/mgit/Mort/timer
BenchmarkTimerSleep-4   	    2000	   1216306 ns/op	       0 B/op	       0 allocs/op
BenchmarkTimeSleep-4    	    1000	   1240398 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	hslam.com/mgit/Mort/timer	3.927s*/
