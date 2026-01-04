package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mirkobrombin/go-revert/v2/pkg/workflow"
)

func main() {
	wf := workflow.New()

	wf.Add("Flaky API",
		workflow.WithRetry(workflow.RetryPolicy{MaxAttempts: 3, Delay: 100 * time.Millisecond},
			func(ctx context.Context) error {
				fmt.Println("Trying API...")
				return nil
			},
		),
		nil,
	)

	group := workflow.Group{
		workflow.Step{
			Name: "Upload Image",
			Do: func(ctx context.Context) error {
				fmt.Println("Uploading Image...")
				time.Sleep(500 * time.Millisecond)
				return nil
			},
			Compensate: func(ctx context.Context) error {
				fmt.Println("Undo: Delete Image")
				return nil
			},
		},
		workflow.Step{
			Name: "Upload Thumbnail",
			Do: func(ctx context.Context) error {
				fmt.Println("Uploading Thumb...")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				fmt.Println("Undo: Delete Thumb")
				return nil
			},
		},
	}
	wf.AddGroup(group)

	fmt.Println("--- Starting Advanced Saga ---")
	if err := wf.Run(context.Background()); err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Println("--- Success ---")
	}
}
