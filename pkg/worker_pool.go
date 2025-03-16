package pkg

import (
	"bad-code-review/model"
	"context"
	"fmt"
	"sync"
)

type UserDetailsWorkerPool struct {
	numOfWorkers int
}

func NewUserDetailsWorkerPool(numWorkersInt int) UserDetailsWorkerPool {
	return UserDetailsWorkerPool{
		numOfWorkers: numWorkersInt,
	}
}

func (w *UserDetailsWorkerPool) ProcessOrdersIds(ctx context.Context, orderIDs []string, jobFunc func(context.Context, string) (model.User, error)) ([]model.User, error) {
	jobsCh, outCh := make(chan string, len(orderIDs)), make(chan model.User)
	wg := &sync.WaitGroup{}
	for i := 0; i < w.numOfWorkers; i++ {
		wg.Add(1)
		go w.processJob(ctx, i, jobsCh, outCh, wg, jobFunc)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()

	go func() {
		for _, orderId := range orderIDs {
			jobsCh <- orderId
		}
		close(jobsCh)
	}()

	users := make([]model.User, 0, len(orderIDs))
	for user := range outCh {
		users = append(users, user)
	}
	return users, nil
}

func (w *UserDetailsWorkerPool) processJob(
	ctx context.Context,
	workerId int,
	jobsCh <-chan string,
	outCh chan<- model.User,
	wg *sync.WaitGroup,
	jobFunc func(context.Context, string) (model.User, error),
) {
	defer wg.Done()
	fmt.Println(fmt.Sprintf("Worker %d started", workerId))
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println(fmt.Sprintf("Worker %d done since context was cancelled", workerId))
			return
		}
	}()

	for jobId := range jobsCh {
		userData, err := jobFunc(ctx, jobId)
		if err != nil {
			fmt.Println(fmt.Sprintf("Worker %d failed to fetch user details for job %s", workerId, jobId))
			continue
		}
		outCh <- userData
	}
}
