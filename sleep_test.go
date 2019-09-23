package timer

import (
	"testing"
	"time"
)

func BenchmarkTimerSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sleep(time.Millisecond*1000)
	}
}

func BenchmarkTimeSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond*1000)
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
