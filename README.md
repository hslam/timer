# timer
A timer library reduced %CPU in top from 20-25% to 2-3% during the idle period of system.

By writing a cgo function Sleep that calls C's usleep directly.

## Feature
Much lower CPU overhead when using timer.Ticker or timer.Sleep()

System |func   |1 us   |100 us|1 ms
 ---- | ----- | ------  | ------  | ------
go     |Ticker     |146.3% |54.0% |12%
go     |Sleep      |130.7% |24.2% |7.9%
cgo    |Ticker     |124.8% |22.5% |6.6%
cgo    |FuncTicker |46.5%  |5.2%  |1.9%
cgo    |Sleep      |45.6%  |4.6%  |1.6%

## Get started

### Install
```
go get hslam.com/mgit/Mort/timer
```
### Import
```
import "hslam.com/mgit/Mort/timer"
```
### Usage
#### Example
```
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
```


### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
timer was written by Mort Huang.


