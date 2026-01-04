package step

import "context"

// Basic is a helper to create a step from two functions.
type Basic struct {
	Name      string
	OnExecute func(ctx context.Context) error
	OnUndo    func(ctx context.Context) error
}

func (b *Basic) Execute(ctx context.Context) error {
	if b.OnExecute != nil {
		return b.OnExecute(ctx)
	}
	return nil
}

func (b *Basic) Compensate(ctx context.Context) error {
	if b.OnUndo != nil {
		return b.OnUndo(ctx)
	}
	return nil
}

// Ensure Basic implements Step
var _ Step = (*Basic)(nil)
