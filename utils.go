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
