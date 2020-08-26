// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import "time"

func min(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

func max(a, b time.Duration) time.Duration {
	if a >= b {
		return a
	}
	return b
}
