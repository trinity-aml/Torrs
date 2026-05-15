package utils

import (
	"sync"
	"sync/atomic"
)

func ParallelFor(begin, end, lim int, fn func(i int)) {
	if end <= begin {
		return
	}
	var wg sync.WaitGroup
	wg.Add(end - begin)
	limits := make(chan struct{}, lim)
	for i := begin; i < end; i++ {
		limits <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-limits }()
			fn(i)
		}(i)
	}
	wg.Wait()
}

func ParallelForEnd(begin, end, lim int, fn func(i int) bool) {
	if end <= begin {
		return
	}
	var wg sync.WaitGroup
	wg.Add(end - begin)
	limits := make(chan struct{}, lim)
	var isEnd atomic.Bool
	for i := begin; i < end; i++ {
		limits <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-limits }()
			if isEnd.Load() {
				return
			}
			if !fn(i) {
				isEnd.Store(true)
			}
		}(i)
	}
	wg.Wait()
}
