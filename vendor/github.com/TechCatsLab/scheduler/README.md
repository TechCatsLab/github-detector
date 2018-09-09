# scheduler

> Warning: Don't use `Schedule` in a Task. It often leads to deadlock. For an example:

```go
type RetryTask struct {
    task Task
    ctx *context.Context
}

func (rt *RetryTask) Do(ctx *context.Context) error {
    if err := rt.task.Do(rt.ctx); err == nil {
        return nil
    }

    // get times and pool from context
    if times > 0 {
        times--
        pool.Schedule(NewRetryContext(times, pool), NewRetryTask(info, rt.task))
        return nil
    }

    return err
}
```
