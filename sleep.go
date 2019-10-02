// +build !use_cgo

package timer

import (
	"time"
)
const (
	ALPHA = 0.118
	BETA  = 2
)
var  Tag = "!use_cgo"

func Sleep(d time.Duration) {
	time.Sleep(d)
}
