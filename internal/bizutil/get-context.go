package bizutil

import (
	"context"
	"time"
)

// GetContext is to get context based on the passing excution timeout
func GetContext(executionTimeout int) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if executionTimeout > 0 {
		return context.WithTimeout(ctx, time.Duration(executionTimeout)*time.Second)
	}
	return context.WithCancel(ctx)
}
