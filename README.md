# timer
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/timer)](https://pkg.go.dev/github.com/hslam/timer)
[![Build Status](https://github.com/hslam/timer/workflows/build/badge.svg)](https://github.com/hslam/timer/actions)
[![codecov](https://codecov.io/gh/hslam/timer/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/timer)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/timer)](https://goreportcard.com/report/github.com/hslam/timer)
[![LICENSE](https://img.shields.io/github/license/hslam/timer.svg?style=flat-square)](https://github.com/hslam/timer/blob/master/LICENSE)

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
```go
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


### Author
timer was written by Meng Huang.


