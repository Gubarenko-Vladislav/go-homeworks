package cli

import (
	"context"
	"fmt"
	"homework-2/internal/workerpool"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func RootCmd(ctx context.Context) *cobra.Command {
	var from int
	var to int
	var workers int
	var timeout time.Duration

	rootCmd := &cobra.Command{
		Use:   "movie-loader",
		Short: "Download movie descriptions",
		Long:  `CLI application for downloading movie descriptions from homeworksite.site.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if from > to {
				return fmt.Errorf("from must be less than or equal to to")
			}

			if workers <= 0 {
				return fmt.Errorf("workers must be greater than 0")
			}

			if timeout <= 0 {
				return fmt.Errorf("timeout must be greater than 0")
			}

			jobs := make(chan int, workers)
			results := make(chan workerpool.Result, workers)

			var wg sync.WaitGroup

			for i := 1; i <= workers; i++ {
				wg.Add(1)
				go workerpool.Worker(timeout, jobs, results, &wg, ctx)
			}

			go func() {
				defer close(jobs)

				for j := from; j <= to; j++ {
					select {
					case jobs <- j:
					case <-ctx.Done():
						return
					}
				}
			}()

			go func() {
				wg.Wait()
				close(results)
			}()

			for res := range results {
				if res.Err != nil {
					fmt.Println(res.Err)
					continue
				}

				fmt.Printf(
					"%d — %s— %d — %s\n",
					res.Res.ID,
					res.Res.Title,
					res.Res.Year,
					res.Res.Director,
				)
			}

			return nil
		},
	}

	rootCmd.Flags().IntVar(
		&from,
		"from",
		0,
		"ID of the first movie",
	)

	rootCmd.Flags().IntVar(
		&to,
		"to",
		0,
		"ID of the last movie",
	)

	rootCmd.Flags().IntVar(
		&workers,
		"workers",
		10,
		"Number of workers",
	)

	rootCmd.Flags().DurationVar(
		&timeout,
		"timeout",
		5*time.Second,
		"HTTP request timeout",
	)

	if err := rootCmd.MarkFlagRequired("from"); err != nil {
		panic(err)
	}

	if err := rootCmd.MarkFlagRequired("to"); err != nil {
		panic(err)
	}

	return rootCmd
}
