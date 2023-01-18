package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var task1 Task = func(ctx context.Context) TaskResult {
	time.Sleep(time.Second * 2)
	return TaskResult{
		"Task1", nil,
	}
}

var task2 Task = func(ctx context.Context) TaskResult {
	time.Sleep(time.Second * 1)
	return TaskResult{
		"Task2", nil,
	}
}

var task3 Task = func(ctx context.Context) TaskResult {
	time.Sleep(time.Second * 3)
	return TaskResult{
		"Task3", nil,
	}
}

// Tests for AnyOneReturn
func TestHappyCaseOfAnyReturn(t *testing.T) {
	startTime := time.Now()
	retStub := AsyncExecutorForFirstReturn(context.TODO(), 3500, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	if retStub.GetResultWhenAnyTaskReturns().Err != nil {
		t.Error("Unexpected error occurred.")
	}
	str, _ := retStub.GetResultWhenAnyTaskReturns().Result.(string)
	if str != "Task2" {
		t.Error("unexpected return value.")
		return
	}
	if time.Since(startTime).Milliseconds() > 1050 {
		t.Error("Time spent is unexpected.")
		return
	}
	fmt.Println(retStub.GetResultWhenAnyTaskReturns().Result.(string),
		retStub.GetResultWhenAnyTaskReturns().Err)
}

func TestTimeoutCaseOfAnyReturn(t *testing.T) {
	startTime := time.Now()
	retStub := AsyncExecutorForFirstReturn(context.TODO(), 100, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	err := retStub.GetResultWhenAnyTaskReturns().Err
	if err == nil || err != ErrTimeout {
		t.Error("Timeout error is expected.")
		return
	}
	if time.Since(startTime).Milliseconds() > 110 {
		t.Error("Time spent is unexpected.")
		return
	}
}

func TestCancelCaseOfAnyReturn(t *testing.T) {
	startTime := time.Now()
	ctx, cacelFn := context.WithCancel(context.Background())
	retStub := AsyncExecutorForFirstReturn(ctx, 100, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	cacelFn()
	err := retStub.GetResultWhenAnyTaskReturns().Err
	if err == nil || err != ErrCancelled {
		t.Error("Cancelled error is expected.")
		return
	}
	if time.Since(startTime).Milliseconds() > 3 {
		t.Error("Time spent is unexpected.")
		return
	}
	fmt.Println(retStub.GetResultWhenAnyTaskReturns())
}

// Tests for AllReturn
func TestHappyCaseOfAllReturn(t *testing.T) {
	startTime := time.Now()
	retStub := AsyncExecutorForAllReturn(context.TODO(), 3500, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	if retStub.GetResultsWhenAllTasksReturn().Err != nil {
		t.Error("Unexpected error occurred.")
		return
	}
	taskRet := *retStub.GetResultsWhenAllTasksReturn().Results
	if taskRet[0] == nil || taskRet[1] == nil || taskRet[2] == nil {
		t.Error("Unexpected results.")
		return
	}
	rets := *retStub.GetResultsWhenAllTasksReturn().Results
	fmt.Println(rets[0], rets[1], rets[2])
	if time.Since(startTime).Milliseconds() > 3050 {
		t.Error("Time spent is unexpected.")
		return
	}
}

func TestTimeoutCaseOfAllReturn(t *testing.T) {
	startTime := time.Now()
	retStub := AsyncExecutorForAllReturn(context.TODO(), 2100, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	rets := *retStub.GetResultsWhenAllTasksReturn().Results
	if rets[0] == nil || rets[1] == nil || rets[2] != nil {
		t.Error("unexpected results", rets[0], rets[1], rets[2])
		return
	}
	fmt.Println(rets[0], rets[1], rets[2])
	err := retStub.GetResultsWhenAllTasksReturn().Err
	fmt.Println(err)
	if err == nil || err != ErrTimeout {
		t.Error("Timeout error is expected.")
		return
	}
	if time.Since(startTime).Milliseconds() > 2200 {
		t.Error("Time spent is unexpected.")
		return
	}
}

func TestCancelCaseOfAllReturn(t *testing.T) {
	startTime := time.Now()
	ctx, cacelFn := context.WithCancel(context.Background())
	retStub := AsyncExecutorForAllReturn(ctx, 100, task1, task2, task3)

	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	cacelFn()
	err := retStub.GetResultsWhenAllTasksReturn().Err
	if err == nil || err != ErrCancelled {
		t.Error("Cancelled error is expected.")
		return
	}
	if time.Since(startTime).Milliseconds() > 3 {
		t.Error("Time spent is unexpected.")
		return
	}
	fmt.Println(retStub.GetResultsWhenAllTasksReturn())
	rets := *retStub.GetResultsWhenAllTasksReturn().Results
	if rets[0] != nil || rets[1] != nil || rets[2] != nil {
		t.Error("unexpected results", rets[0], rets[1], rets[2])
		return
	}
	fmt.Println(rets[0], rets[1], rets[2])
}
