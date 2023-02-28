package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

var testTask = func(ctx context.Context) TaskResult {
	time.Sleep(time.Millisecond * 300)
	return TaskResult{Result: "OK", Err: nil, CostMs: 200}
}
var hystrixName string

func BuildHystrix(timeout int) {
	hystrixName = "async_hystrix"
	hystrix.ConfigureCommand(hystrixName, hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  2,
		RequestVolumeThreshold: 2,
		SleepWindow:            timeout * 2,
		ErrorPercentThreshold:  5,
	})
}
func TestHystrixHappyPath(t *testing.T) {
	BuildHystrix(1000)
	startTime := time.Now()

	retStub := AsynHystrix(context.TODO(), testTask, hystrixName, 1000)
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

func TestHystrixTimeoutPath(t *testing.T) {
	BuildHystrix(100)
	startTime := time.Now()
	retStub := AsynHystrix(context.TODO(), testTask, hystrixName, 100)
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

	// test hystrix timeout
	retStub = AsynHystrix(context.TODO(), testTask, hystrixName, 1000)
	ret = retStub.GetResult()
	if ret.Err == nil {
		t.Error("It is not unexpected execution time.")
	}
	if !errors.Is(ret.Err, hystrix.ErrTimeout) {
		t.Error("It is unexpected error", ret.Err)
	}
}

func TestHystrixConcurrencyPath(t *testing.T) {
	BuildHystrix(1000)
	// test hystrix max cocurrency
	stubs := []ResultStub{}
	stubs = append(stubs, AsynHystrix(context.TODO(), testTask, hystrixName, 1000))
	time.Sleep(time.Millisecond * 100)
	stubs = append(stubs, AsynHystrix(context.TODO(), testTask, hystrixName, 1000))
	time.Sleep(time.Millisecond * 100)
	stubs = append(stubs, AsynHystrix(context.TODO(), testTask, hystrixName, 1000))

	ret := stubs[0].GetResult()
	if ret.Err != nil {
		t.Error("It is unexpected execution time.", ret.Err)
	}

	ret = stubs[1].GetResult()
	if ret.Err != nil {
		t.Error("It is unexpected execution time.", ret.Err)
	}

	ret = stubs[2].GetResult()
	if ret.Err == nil {
		t.Error("It is not unexpected execution time.")
	}
	if !errors.Is(ret.Err, hystrix.ErrMaxConcurrency) {
		t.Error("It is unexpected error", ret.Err)
	}
}
