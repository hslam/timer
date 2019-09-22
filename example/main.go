package main

import (
	"flag"
	"time"
	"hslam.com/mgit/Mort/timer"
	"fmt"
)
func main(){
	t := *flag.String("t", "funcTicker", "use timer")
	flag.Parse()
	switch t {
	case "ticker":
		ticker:=timer.NewTicker(time. Millisecond)
		defer ticker.Stop()
		for range ticker.C {
		}
	case "funcTicker":
		ticker:=timer.NewFuncTicker(time.Millisecond,nil)
		defer ticker.Stop()
		ticker.Tick(func() {
			//todo
		})
		select {
		}
	case "sleep":
		for{
			timer.Sleep(time.Millisecond)
		}
	default:
		fmt.Println("use  ticker, funcTicker or sleep")
	}
}