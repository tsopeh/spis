package main

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createLimitedWaitGroup(maxParallelJobsCount int) (func(func()), func()) {
	c := make(chan struct{}, maxParallelJobsCount)

	runJob := func(job func()) {
		c <- struct{}{}
		go func() {
			job()
			<-c
		}()
	}

	wait := func() {
		close(c)
	}

	return runJob, wait
}
