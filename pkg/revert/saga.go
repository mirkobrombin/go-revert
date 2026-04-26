package revert

import (
	"github.com/mirkobrombin/go-revert/v2/pkg/workflow"
)

type Workflow = workflow.Workflow
type Step = workflow.Step
type Group = workflow.Group
type RetryPolicy = workflow.RetryPolicy

var New = workflow.New
var WithRetry = workflow.WithRetry
