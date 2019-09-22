package main

import (
	"flag"
	"time"
	"hslam.com/mgit/Mort/timer"
	"fmt"
)
func main(){
	t := *flag.String("t", "sleep", "use timer")
	flag.Parse()
	switch t {
	case "ticker":
		ticker:=timer.NewTicker(time. Millisecond)
		defer ticker.Stop()
		for range ticker.C {
		}
	case "tickFunc":
		ticker:=timer.NewTicker(time.Millisecond)
		defer ticker.Stop()
		ticker.TickFunc(func() {
			//todo
		})
		select {
		}
	case "sleep":
		for{
			timer.Sleep(time.Microsecond*100)
		}
	default:
		fmt.Println("use  ticker, tickFunc or sleep")
	}
}