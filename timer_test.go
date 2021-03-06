// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package timer

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimerSleep(t *testing.T) {
	num := 1000
	d := time.Millisecond
	run := d * time.Duration(num)

	startTime := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		time.Sleep(d)
	}
	runTime := time.Now().UnixNano() - startTime
	ratioTime := float64(runTime-int64(run)) / float64(run)

	startTimer := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		Sleep(d)
	}
	runTimer := time.Now().UnixNano() - startTimer
	ratioTimer := float64(runTimer-int64(run)) / float64(run)

	if ratioTimer > ratioTime && ratioTimer-ratioTime > 0.1 {
		t.Errorf("error ratio Timer:%.2f Time:%.2f", ratioTimer, ratioTime)
	}
	t.Log("timer.Sleep()", runTimer, ",time.Sleep()", runTime, ",ratioTimer-ratioTime", ratioTimer-ratioTime)
}

func BenchmarkTimerSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sleep(time.Millisecond)
	}
}

func TestTimerTicker(t *testing.T) {
	num := 1000
	d := time.Millisecond
	timeTicker := time.NewTicker(d)
	run := d * time.Duration(num)

	startTime := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		<-timeTicker.C
	}
	runTime := time.Now().UnixNano() - startTime
	timeTicker.Stop()
	ratioTime := float64(runTime-int64(run)) / float64(run)

	timerTicker := NewTicker(d)
	startTimer := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		<-timerTicker.C
	}
	runTimer := time.Now().UnixNano() - startTimer
	timerTicker.Stop()
	ratioTimer := float64(runTimer-int64(run)) / float64(run)

	if math.Abs(ratioTimer-ratioTime) > 0.1 || math.Abs(ratioTimer) > 0.1 {
		t.Errorf("error ratio Timer:%.2f Time:%.2f", ratioTimer, ratioTime)
	}
	t.Log("timer.Ticker", runTimer, ",time.Ticker", runTime, ",ratioTimer", ratioTimer, ",ratioTime", ratioTime)
	timerTicker.Stop()
}

func TestNewTicker(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error()
		}
	}()
	NewTicker(0)
}

func TestF(t *testing.T) {
	c := make(chan time.Time, 1)
	f(c)
	f(c)
}

func TestTick(t *testing.T) {
	if Tick(0) != nil {
		t.Error()
	}
	when := <-Tick(time.Millisecond)
	if when.After(time.Now()) {
		t.Error()
	}
}

func TestTickFunc(t *testing.T) {
	once := int32(0)
	ch := make(chan struct{}, 1)
	ticker := TickFunc(time.Millisecond, func() {
		if atomic.CompareAndSwapInt32(&once, 0, 1) {
			close(ch)
		}
	})
	time.Sleep(time.Second)
	<-ch
	ticker.Stop()

}

func TestTickFuncMore(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error()
		}
	}()
	TickFunc(0, func() {

	})
}

func TestRecoverTimer(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error()
			}
		}()
		NewTimer(0)
	}()
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error()
			}
		}()
		NewTimer(time.Millisecond).Reset(0)
	}()
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error()
			}
		}()
		AfterFunc(0, func() {
		})
	}()
}

func TestTimer(t *testing.T) {
	{
		when := <-NewTimer(time.Millisecond).C
		if when.After(time.Now()) {
			t.Error()
		}
	}
	{
		when := <-After(time.Millisecond)
		if when.After(time.Now()) {
			t.Error()
		}
	}
	{
		timer := NewTimer(time.Millisecond)
		timer.Reset(time.Millisecond)
		timer.Stop()
		timer.Stop()
	}
	{
		once := int32(0)
		ch := make(chan struct{}, 1)
		timer := AfterFunc(time.Millisecond, func() {
			if atomic.CompareAndSwapInt32(&once, 0, 1) {
				close(ch)
			}
		})
		<-ch
		timer.r.Stop()
		timer.Stop()
	}
	end := time.Now().Add(time.Second * 6)
	wg := sync.WaitGroup{}
	for i := time.Duration(0); time.Now().Before(end); i++ {
		wg.Add(1)
		go func(v time.Duration) {
			d := time.Now().Sub(end)
			if -d-v > time.Microsecond {
				<-NewTimer(-d - v).C
			}
			wg.Done()
		}(i)
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
	time.Sleep(idleTime * 2)
}

func TestTrigger(t *testing.T) {
	c := make(chan struct{}, 1)
	trigger(c)
	trigger(c)
	resetTrigger(c)
	resetTriggerOnce(c)
}

func TestWhen(t *testing.T) {
	if when(0) < 0 {
		t.Error()
	}
	if when(1<<63-1) < 0 {
		t.Error()
	}
}
