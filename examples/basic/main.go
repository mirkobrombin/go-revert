package main

import (
	"context"
	"fmt"

	"github.com/mirkobrombin/go-revert/v2/pkg/engine"
	"github.com/mirkobrombin/go-revert/v2/pkg/step"
)

func main() {
	wf := engine.New()

	// Step 1: Reserve Stock
	wf.Add(&step.Basic{
		Name: "Reserve Stock",
		OnExecute: func(ctx context.Context) error {
			fmt.Println("Step 1: Reserving Stock...")
			return nil
		},
		OnUndo: func(ctx context.Context) error {
			fmt.Println("Undo 1: Releasing Stock")
			return nil
		},
	})

	// Step 2: Charge Card (Fails)
	wf.Add(&step.Basic{
		Name: "Charge Card",
		OnExecute: func(ctx context.Context) error {
			fmt.Println("Step 2: Charging Card...")
			return fmt.Errorf("insufficient funds")
		},
		OnUndo: func(ctx context.Context) error {
			fmt.Println("Undo 2: Refund Card")
			return nil
		},
	})

	fmt.Println("Starting Workflow...")
	if err := wf.Run(context.Background()); err != nil {
		fmt.Printf("Workflow failed: %v\n", err)
	} else {
		fmt.Println("Workflow completed successfully!")
	}
}
