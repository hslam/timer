// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !use_cgo

package timer

import (
	"time"
)

var Tag = "!use_cgo"

func Sleep(d time.Duration) {
	time.Sleep(d)
}
