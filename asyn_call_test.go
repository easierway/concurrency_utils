package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

var task = func(ctx context.Context) TaskResult {
	time.Sleep(time.Second * 3)
	return TaskResult{"OK", nil}
}

func TestHappyPath(t *testing.T) {
	startTime := time.Now()
	retStub := AsynExecutor(context.TODO(), task, 5000)
	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	ret := retStub.GetResult()
	if ret.Err != nil {
		t.Error(ret.Err)
	}
	if time.Since(startTime).Microseconds() < 3000 {
		t.Error("It is not unexpected execution time.")
	}
	fmt.Println(ret.Result.(string), ret.Err)
}

func TestTimeoutPath(t *testing.T) {
	startTime := time.Now()
	retStub := AsynExecutor(context.TODO(), task, 100)
	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	ret := retStub.GetResult()
	if ret.Err == nil {
		t.Error("It is not unexpected execution time.")
	}
	if !errors.Is(ret.Err, ErrTimeout) {
		t.Error("It is unexpected error", ret.Err)
	}
	fmt.Println(ret.Result, ret.Err)
}

func TestCancelPath(t *testing.T) {
	startTime := time.Now()
	ctx, cancelFn := context.WithCancel(context.Background())
	retStub := AsynExecutor(ctx, task, 1)
	if time.Since(startTime).Milliseconds() > 2 {
		t.Error("It is not asynchronous call.")
	}
	fmt.Println("nonblocking")
	cancelFn()
	ret := retStub.GetResult()
	if ret.Err == nil {
		t.Error("It is not unexpected execution time.")
	}
	if !errors.Is(ret.Err, ErrCancelled) {
		t.Error("It is unexpected error", ret.Err)
	}
	fmt.Println(ret.Result, ret.Err)
}
