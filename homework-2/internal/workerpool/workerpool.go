package workerpool

import (
	"context"
	"homework-2/internal/clieant"
	"homework-2/internal/movie"
	"sync"
	"time"
)

type Result struct {
	Res movie.Movie
	Err error
}

func Worker(
	timeout time.Duration,
	jobs <-chan int,
	results chan<- Result,
	wg *sync.WaitGroup,
	ctx context.Context,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			res, err := clieant.Get(job, timeout, ctx)

			select {
			case results <- Result{Res: res, Err: err}:
			case <-ctx.Done():
				return
			}
		}
	}
}
