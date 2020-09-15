// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !use_cgo

package timer

import (
	"time"
)

// Tag is the sleep type.
var Tag = "!use_cgo"

// Sleep pauses the current goroutine for at least the duration d.
// A negative or zero duration causes Sleep to return immediately.
func Sleep(d time.Duration) {
	time.Sleep(d)
}
