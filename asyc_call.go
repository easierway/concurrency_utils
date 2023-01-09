package concurrency

import (
	"context"
	"errors"
	"time"
)

var ErrTimeout = errors.New("time out")
var ErrCancelled = errors.New("task has been cancelled")

type TaskResult struct {
	Result interface{}
	Err    error
}

type ResultStub struct {
	retCh     chan TaskResult
	timeoutMs int64
	ctx       context.Context
}

func (rs *ResultStub) GetResult() TaskResult {
	var tResult TaskResult
	timer := time.NewTimer(time.Millisecond * time.Duration(rs.timeoutMs))
	select {
	case ret := <-rs.retCh:
		tResult.Result = ret
		return tResult
	case <-timer.C:
		tResult.Err = ErrTimeout
		return tResult
	case <-rs.ctx.Done():
		tResult.Err = ErrCancelled
		return tResult
	}
	return tResult
}

type Task func(ctx context.Context) TaskResult

// AsynExecutor is to run your text in asynchronous way,
// task is the task need to be executed
// timeoutMs is the timeout setting for the task execution
// ResultStub ResultStub.GetResult() is to get the result of the task.
func AsynExecutor(ctx context.Context, task Task,
	timeoutMs int64) *ResultStub {
	retCh := make(chan TaskResult, 1)
	go func() {
		retCh <- task(ctx)
	}()
	return &ResultStub{retCh, timeoutMs, ctx}
}
