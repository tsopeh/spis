package main

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
	return SemaphoredWaitGroup{
		AddJob: func(job func()) {
			c <- struct{}{}
			go func() {
				job()
				<-c
			}()
		},
		WaitAll: func() {
			close(c)
		},
	}
}
