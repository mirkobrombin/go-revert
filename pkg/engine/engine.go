package engine

import (
	"context"
	"fmt"

	"github.com/mirkobrombin/go-revert/v2/pkg/step"
)

// Workflow manages a sequence of steps and handles compensations on failure.
type Workflow struct {
	steps    []step.Step
	executed []step.Step
}

// New creates a new Workflow.
func New() *Workflow {
	return &Workflow{
		steps:    make([]step.Step, 0),
		executed: make([]step.Step, 0),
	}
}

// Add appends a step to the workflow.
func (w *Workflow) Add(s step.Step) *Workflow {
	w.steps = append(w.steps, s)
	return w
}

// Run executes the workflow steps in order.
// If a step fails, it executes compensations for all previously successful steps in reverse order.
func (w *Workflow) Run(ctx context.Context) error {
	for _, s := range w.steps {
		select {
		case <-ctx.Done():
			w.compensate(ctx)
			return ctx.Err()
		default:
			if err := s.Execute(ctx); err != nil {
				w.compensate(ctx)
				return fmt.Errorf("step execution failed: %w", err)
			}
			w.executed = append(w.executed, s)
		}
	}
	return nil
}

func (w *Workflow) compensate(ctx context.Context) {
	// Root context might be cancelled, use a new one for compensation if needed
	// but usually we respect the original ctx for timeouts if they apply to the whole flow.
	for i := len(w.executed) - 1; i >= 0; i-- {
		s := w.executed[i]
		if err := s.Compensate(ctx); err != nil {
			// In a real scenario, we might want to log this or retry.
			fmt.Printf("warning: compensation failed for step: %v\n", err)
		}
	}
}

// Register (Declarative) allows registering multiple steps from a struct
// (Future expansion: discover steps via reflection or tags)
func (w *Workflow) Register(steps ...step.Step) *Workflow {
	for _, s := range steps {
		w.Add(s)
	}
	return w
}
