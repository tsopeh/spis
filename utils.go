package main

import "sync"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type SemaphoredWaitGroup struct {
	AddJob  func(func())
	WaitAll func()
}

func CreateSemaphoredWaitGroup(maxParallelJobsCount int) SemaphoredWaitGroup {
	c := make(chan struct{}, maxParallelJobsCount)
	wg := sync.WaitGroup{}
	return SemaphoredWaitGroup{
		AddJob: func(job func()) {
			wg.Add(1)
			c <- struct{}{}
			go func() {
				job()
				<-c
				wg.Done()
			}()
		},
		WaitAll: func() {
			wg.Wait()
			close(c)
		},
	}
}
