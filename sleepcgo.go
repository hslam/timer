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
	ALPHA = 0.1
	BETA  = 2
)
var  Tag = "use_cgo"

func Sleep(d time.Duration) {
	if d < time.Microsecond {
		d=time.Microsecond
	}
	var duration C.uint
	duration=C.uint(int64(d)/1000)
	C.usleep(duration)
}
