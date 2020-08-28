package main

import (
	"flag"
	"fmt"
	"github.com/hslam/timer"
)

func main() {
	t := *flag.String("t", "Sleep", "use timer")
	flag.Parse()
	fmt.Println(timer.Tag)
	switch t {
	case "Ticker":
		t := timer.NewTicker(timer.Millisecond)
		defer t.Stop()
		for range t.C {
			//todo
		}
	case "TickFunc":
		t := timer.TickFunc(timer.Millisecond, func() {
			//todo
		})
		defer t.Stop()
		select {}
	case "Timer":
		t := timer.NewTimer(timer.Millisecond)
		defer t.Stop()
		for range t.C {
			t.Reset(timer.Millisecond)
			//todo
		}
	case "Sleep":
		for {
			timer.Sleep(timer.Millisecond)
			//todo
		}
	case "After":
		for {
			select {
			case <-timer.After(timer.Millisecond):
				//todo
			}
		}
	default:
		fmt.Println("use Ticker, TickFunc, Timer, Sleep or After")
	}
}
