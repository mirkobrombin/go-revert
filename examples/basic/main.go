package main

import (
	"context"
	"fmt"

	"github.com/mirkobrombin/go-revert/v2/pkg/workflow"
)

func main() {
	wf := workflow.New()

	wf.Add("Reserve Stock",
		func(ctx context.Context) error {
			fmt.Println("Step 1: Reserving Stock...")
			return nil
		},
		func(ctx context.Context) error {
			fmt.Println("Undo 1: Releasing Stock")
			return nil
		},
	)

	wf.Add("Charge Card",
		func(ctx context.Context) error {
			fmt.Println("Step 2: Charging Card...")
			return fmt.Errorf("insufficient funds")
		},
		func(ctx context.Context) error {
			fmt.Println("Undo 2: Refund Card")
			return nil
		},
	)

	fmt.Println("Starting Workflow...")
	if err := wf.Run(context.Background()); err != nil {
		fmt.Printf("Workflow failed: %v\n", err)
	} else {
		fmt.Println("Workflow completed successfully!")
	}
}
