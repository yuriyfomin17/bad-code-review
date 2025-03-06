package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func WorkerPoolWithContext() {
	numOfWorkers := 3
	jobsCh, outCh := make(chan string), make(chan string)
	jobsIds := []string{"1", "2", "3", "4", "5", "6"}
	wg := &sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go func() {
			timeout, _ := context.WithTimeout(context.Background(), time.Second*3)
			processJob(timeout, i, jobsCh, outCh, wg)
		}()
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	go func() {
		for _, jobId := range jobsIds {
			jobsCh <- jobId
		}
		close(jobsCh)
	}()

	for val := range outCh {
		fmt.Println(fmt.Sprintf("Worker %s done", val))
	}
}

func processJob(ctx context.Context, workerId int, jobsCh <-chan string, outCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(fmt.Sprintf("Worker %d started", workerId))
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println(fmt.Sprintf("Worker %d done since timeout was reached", workerId))
				return
			}
		}
	}()
	for jobId := range jobsCh {
		fmt.Println(fmt.Sprintf("Worker %d processing job %s", workerId, jobId))
		time.Sleep(time.Second)
		outCh <- jobId
	}
}
