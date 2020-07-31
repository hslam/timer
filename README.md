# timer
A timer library reduced %CPU in top from 20-25% to 2-3% during the idle period of system.

By writing a cgo function Sleep that calls C's usleep directly.

Much lower CPU overhead when using timer.Ticker or timer.Sleep()

System |func   |1 us   |100 us|1 ms
 ---- | ----- | ------  | ------  | ------
go     |Ticker     |146.3% |54.0% |12%
go     |Sleep      |130.7% |24.2% |7.9%
cgo    |Ticker     |124.8% |22.5% |6.6%
cgo    |FuncTicker |46.5%  |5.2%  |1.9%
cgo    |Sleep      |45.6%  |4.6%  |1.6%

## Feature
* cgo
* Ticker
* FuncTicker
* Timer
* Sleep
* After

## Get started

### Install
```
go get github.com/hslam/timer
```
### Import
```
import "github.com/hslam/timer"
```
### Usage
#### Example
```
package main
import (
	"flag"
	"github.com/hslam/timer"
	"fmt"
)
func main(){
	t := *flag.String("t", "Sleep", "use timer")
	flag.Parse()
	fmt.Println(timer.Tag)
	switch t {
	case "Ticker":
		t:=timer.NewTicker(timer.Millisecond)
		defer t.Stop()
		for range t.C {
			//todo
		}
	case "FuncTicker":
		t:=timer.NewFuncTicker(timer.Millisecond,nil)
		defer t.Stop()
		t.Tick(func() {
			//todo
		})
		select {}
	case "Timer":
		t:=timer.NewTimer(timer.Millisecond)
		defer t.Stop()
		for range t.C {
			t.Reset(timer.Millisecond)
			//todo
		}
	case "Sleep":
		for{
			timer.Sleep(timer.Millisecond)
			//todo
		}
	case "After":
		for{
			select {
			case <-timer.After(timer.Millisecond):
				//todo
			}
		}
	default:
		fmt.Println("use Ticker,FuncTicker,Timer,Sleep or After")
	}
}
```

### Build
```
go build -tags=use_cgo
```

### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Authors
timer was written by Meng Huang.


