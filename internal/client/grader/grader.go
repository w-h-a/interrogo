package grader

import (
	"context"

	"github.com/w-h-a/interrogo/api/test_result/v1alpha1"
)

type Grader interface {
	Grade(ctx context.Context, history []v1alpha1.Message, policy string) (v1alpha1.TestResult, error)
}
