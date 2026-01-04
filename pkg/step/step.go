package step

import "context"

// Step defines a single reversible operation in a workflow.
type Step interface {
	Execute(ctx context.Context) error
	Compensate(ctx context.Context) error
}

// Metadata holds declarative information about a step.
type Metadata struct {
	Name        string
	Description string
	Critical    bool // If true, failure stops the workflow. Defaults to true.
}
