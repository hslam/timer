// +build !use_cgo

package timer

import (
	"time"
)
const (
	ALPHA = 0.8
	BETA  = 1.3
)

var  Tag = "!use_cgo"

func Sleep(d time.Duration) {
	time.Sleep(d)
}
