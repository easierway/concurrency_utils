package concurrency

import (
	"context"
	"time"
)

// MultiTaskResults is the results when calling GetResultsWhenAllTasksReturn
type MultiTaskResults struct {
	// Results is the results of the tasks
	// could be some of the task results
	Results *[]*TaskResult
	// Err is the errors (including errors about timeout and Cancelled) of the process,
	Err error
}

type AllResultStub struct {
	tTaskResults *[]*TaskResult
	retCh        chan int
	timeoutMs    int64
	ctx          context.Context
	allReturns   *MultiTaskResults
	numOfTasks   int
	timer        *time.Timer
}

type AnyResultStub struct {
	retCh       chan TaskResult
	timeoutMs   int64
	ctx         context.Context
	firstReturn *TaskResult
	timer       *time.Timer
}

func (stub *AnyResultStub) GetResultWhenAnyTaskReturns() TaskResult {
	if stub.firstReturn != nil {
		return *stub.firstReturn
	}
	stub.timer = time.NewTimer(time.Millisecond * time.Duration(stub.timeoutMs))
	select {
	case ret := <-stub.retCh:
		stub.firstReturn = &ret
	case <-stub.timer.C:
		stub.firstReturn = &TaskResult{nil, ErrTimeout, 0}
	case <-stub.ctx.Done():
		stub.firstReturn = &TaskResult{nil, ErrCancelled, 0}
	}
	return *stub.firstReturn

}

func (stub *AllResultStub) GetResultsWhenAllTasksReturn() MultiTaskResults {
	numOfTasks := stub.numOfTasks
	if stub.allReturns != nil {
		return *stub.allReturns
	}
	stub.timer = time.NewTimer(time.Millisecond * time.Duration(stub.timeoutMs))
	stub.allReturns = &MultiTaskResults{
		Results: stub.tTaskResults,
	}
	for i := 0; i < numOfTasks; i++ {
		select {
		case <-stub.retCh:
		case <-stub.timer.C:
			stub.allReturns.Err = ErrTimeout
		case <-stub.ctx.Done():
			stub.allReturns.Err = ErrCancelled
		}
	}
	return *stub.allReturns
}

func AsyncExecutorForFirstReturn(ctx context.Context, timeoutMs int64, tasks ...Task) AnyResultStub {
	numOfTasks := len(tasks)
	resultsCh := make(chan TaskResult, numOfTasks)
	stub := AnyResultStub{
		retCh:     resultsCh,
		timeoutMs: timeoutMs,
		ctx:       ctx,
	}
	for i := 0; i < numOfTasks; i++ {
		go func(id int) {
			ret := tasks[id](ctx)
			resultsCh <- ret

		}(i)
	}
	return stub
}

func AsyncExecutorForAllReturn(ctx context.Context, timeoutMs int64, tasks ...Task) AllResultStub {
	numOfTasks := len(tasks)
	resultsCh := make(chan int, numOfTasks)
	results := make([]*TaskResult, numOfTasks)
	stub := AllResultStub{
		tTaskResults: &results,
		retCh:        resultsCh,
		timeoutMs:    timeoutMs,
		ctx:          ctx,
		numOfTasks:   numOfTasks,
	}
	for i := 0; i < numOfTasks; i++ {
		go func(id int) {
			ret := tasks[id](ctx)
			results[id] = &ret
			resultsCh <- id

		}(i)
	}
	return stub
}
