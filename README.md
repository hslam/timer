# timer
Package timer provides functionality for measuring time. Much lower CPU overhead when using use_cgo tags during the idle period of system.

## Feature
* cgo/go
* Ticker
* TickFunc
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
```

### Build
#### go
```
go build
```
#### cgo
```
go build -tags=use_cgo
```

### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Authors
timer was written by Meng Huang.


