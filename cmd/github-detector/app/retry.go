/*
 * Revision History:
 *     Initial: 2018/08/04        Li Zebang
 */

package app

import (
	"context"
	"errors"

	"github.com/TechCatsLab/scheduler"
)

type (
	// RetryTaskInfo -
	RetryTaskInfo struct {
		Times int
	}

	// RetryTask -
	RetryTask struct {
		ctx  context.Context
		task scheduler.Task
	}
)

// NewRetryTask -
func NewRetryTask(ctx context.Context, task scheduler.Task) *RetryTask {
	return &RetryTask{
		ctx:  ctx,
		task: task,
	}
}

// NewRetryTaskContext -
func NewRetryTaskContext(info *RetryTaskInfo) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, RetryTaskKey, info)

	return ctx
}

// Do -
func (rt *RetryTask) Do(ctx context.Context) (err error) {
	err = rt.task.Do(rt.ctx)
	if err == nil {
		return nil
	}

	info, ok := ctx.Value(RetryTaskKey).(*RetryTaskInfo)
	if !ok {
		return errors.New("assertion fail")
	}

	for index := 0; index < info.Times; index++ {
		err = rt.task.Do(rt.ctx)
		if err == nil {
			return nil
		}
	}

	return err
}
