package concurrency

import (
	"context"
	"errors"
	"time"
)

var ErrTimeout = errors.New("time out")
var ErrCancelled = errors.New("task has been cancelled")

// TaskResult is the result of the asynchronous task
type TaskResult struct {
	Result interface{}
	Err    error
}

// ResultStub is returned immediately after calling AsynExecutor
// It is the stub to get the result.
type ResultStub struct {
	retCh     chan TaskResult
	timeoutMs int64
	ctx       context.Context
}

// GetResult is to get the asynchronous task's result
// it will be blocked util the task is completed/failed.
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

// AsynExecutor is to run your task in asynchronous way,
// task is the task need to be executed
// timeoutMs is the timeout setting for the task execution
// ResultStub ResultStub.GetResult() is to get the result of the task.
func AsynExecutor(ctx context.Context, task Task,
	timeoutMs int64) ResultStub {
	retCh := make(chan TaskResult, 1)
	go func() {
		retCh <- task(ctx)
	}()
	return ResultStub{retCh, timeoutMs, ctx}
}
