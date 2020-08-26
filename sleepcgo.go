// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build use_cgo

package timer

/*
#include <unistd.h>
*/
import "C"
import (
	"time"
)

const (
	ALPHA = 0.8
	BETA  = 1.6
	//defaultMaxSleep	C.uint = 19
)

var Tag = "use_cgo"

func Sleep(d time.Duration) {
	if d < time.Microsecond {
		d = time.Microsecond
	}
	var duration C.uint
	duration = C.uint(int64(d) / 1000)
	C.usleep(duration)
}

//func Sleep(d time.Duration) {
//	if d < time.Microsecond {
//		d=time.Microsecond
//	}
//	var duration C.uint
//	duration=C.uint(int64(d)/1000)
//	if duration>defaultMaxSleep{
//		var sleep C.uint
//		for {
//			C.usleep(defaultMaxSleep)
//			sleep+=defaultMaxSleep
//			if duration-sleep<=defaultMaxSleep{
//				C.usleep(duration-sleep)
//				break
//			}
//		}
//	}else {
//		C.usleep(duration)
//	}
//}
