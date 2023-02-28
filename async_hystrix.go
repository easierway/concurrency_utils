package concurrency

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
)

func AsynHystrix(ctx context.Context, task Task, name string, timeoutMs int64) ResultStub {
	retCh := make(chan TaskResult, 1)

	runC := func(ctx context.Context) error {
		retCh <- task(ctx)
		return nil
	}
	errChan := hystrix.GoC(ctx, name, runC, nil)

	return ResultStub{retCh: retCh,
		timeoutMs: timeoutMs,
		ctx:       ctx,
		errChan:   errChan}
}
