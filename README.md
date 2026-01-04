# Go Revert (Saga Pattern)

**Go Revert** is a minimal, robust library for providing **Application-Level Atomicity** using the Saga Pattern (Forward Compensation).

It allows you to define workflows where each step has a corresponding rollback action. If any step fails (or panics), the library automatically executes the rollback actions of all previously successful steps in reverse order (LIFO), ensuring your system returns to a consistent state.

## Features

- **Atomic Workflows**: Treat heterogenous operations (DB, API, IO) as a single transaction.
- **Panic Safe**: Automatically catches panics in `Do` blocks and triggers rollbacks.
- **Context Aware**: Full support for `context.Context` for cancellation and timeouts.
- **Error Aggregation**: Uses Go 1.20+ `errors.Join` to report both the original failure and any rollback errors.
- **Zero Dependencies**: Pure Go standard library.

## Getting Started

### Installation

```bash
go get github.com/mirkobrombin/go-revert/v2
```

### Basic Usage

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mirkobrombin/go-revert/v2/pkg/revert"
)

func main() {
	wf := revert.New()

	// Create a Resource
	wf.Add(
		func(ctx context.Context) error {
			fmt.Println("Creating file...")
			return nil
		},
		func(ctx context.Context) error {
			fmt.Println("Deleting file (Rollback)")
			return nil
		},
	)

	// Fails!
	wf.Add(
		func(ctx context.Context) error {
			return errors.New("something went wrong")
		},
		func(ctx context.Context) error {
			return nil // Nothing to rollback here usually
		},
	)

	// Execute
	if err := wf.Run(context.Background()); err != nil {
		fmt.Printf("Workflow failed gracefully: %v\n", err)
	}
}
```

## Documentation

- **[Core Concepts](docs/concepts.md)**: How the LIFO stack works.
- **[Error Handling](docs/errors.md)**: Panic recovery and error joining.

## License

MIT License.
