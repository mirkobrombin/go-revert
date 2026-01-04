package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Step struct {
	Name       string
	Do         func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type Group []Step // A Group is just a slice of Steps executed in parallel

type Workflow struct {
	steps []any // Can be Step or Group
	stack []Step
	mu    sync.Mutex
}

func New() *Workflow {
	return &Workflow{
		stack: make([]Step, 0),
	}
}

func (w *Workflow) Add(name string, do, compensate func(ctx context.Context) error) {
	w.steps = append(w.steps, Step{
		Name:       name,
		Do:         do,
		Compensate: compensate,
	})
}

func (w *Workflow) AddGroup(g Group) {
	w.steps = append(w.steps, g)
}

func (w *Workflow) Run(ctx context.Context) error {
	for _, item := range w.steps {
		// Check Context before starting step
		if ctx.Err() != nil {
			return w.rollback(ctx, ctx.Err())
		}

		var err error
		switch v := item.(type) {
		case Step:
			err = w.runStep(ctx, v)
		case Group:
			err = w.runGroup(ctx, v)
		}

		if err != nil {
			return w.rollback(ctx, err)
		}
	}
	return nil
}

func (w *Workflow) runStep(ctx context.Context, step Step) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in step '%s': %v", step.Name, r)
		}
	}()

	if err := step.Do(ctx); err != nil {
		return fmt.Errorf("step '%s' failed: %w", step.Name, err)
	}

	w.mu.Lock()
	if step.Compensate != nil {
		w.stack = append(w.stack, step)
	}
	w.mu.Unlock()

	return nil
}

func (w *Workflow) runGroup(ctx context.Context, group Group) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(group))

	// Temporarily store successful steps in this group to add to stack later
	// If the group fails, we only compensate what succeeded inside the group?
	// Actually, if we use w.runStep() inside goroutine, it appends to w.stack safely.
	// But we must handle partial failure rollback within the group logic or rely on main rollback?
	// For simplicity, we let them append to stack. If one fails, Run() returns error and triggers rollback of everything in stack.

	for _, step := range group {
		wg.Add(1)
		go func(s Step) {
			defer wg.Done()
			if err := w.runStep(ctx, s); err != nil {
				errChan <- err
			}
		}(step)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		var errs []error
		for e := range errChan {
			errs = append(errs, e)
		}
		return errors.Join(errs...)
	}
	return nil
}

func (w *Workflow) rollback(ctx context.Context, triggerErr error) error {
	rollbackCtx := context.WithoutCancel(ctx)
	var errs []error
	errs = append(errs, triggerErr)

	// LIFO
	w.mu.Lock()
	defer w.mu.Unlock()

	for i := len(w.stack) - 1; i >= 0; i-- {
		step := w.stack[i]
		if err := w.safeCompensate(rollbackCtx, step); err != nil {
			errs = append(errs, fmt.Errorf("rollback failed for '%s': %w", step.Name, err))
		}
	}

	return errors.Join(errs...)
}

func (w *Workflow) safeCompensate(ctx context.Context, step Step) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during compensation: %v", r)
		}
	}()
	return step.Compensate(ctx)
}
