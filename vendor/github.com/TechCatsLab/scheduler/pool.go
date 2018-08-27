/*
 * Revision History:
 *     Initial: 2018/07/10        Tong Yuehong
 *     Modify:  2018/08/08        Li Zebang
 */

package scheduler

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

var (
	// ErrScheduleTimeout happens when task schedule failed during the specific interval.
	ErrScheduleTimeout = errors.New("schedule not available currently")

	errSchedulerPoolStop = errors.New("pool stopped")
)

type taskWrapper struct {
	task Task
	ctx  context.Context
}

func (t *taskWrapper) Do(ctx context.Context) error {
	return t.task.Do(ctx)
}

// Pool caches tasks and schedule tasks to work.
type Pool struct {
	queue    chan Task
	workers  chan chan Task
	shutdown chan struct{}
	stop     sync.Once
	wg       sync.WaitGroup
}

// New a goroutine pool.
func New(qsize, wsize int) *Pool {
	if wsize == 0 {
		wsize = runtime.NumCPU()
	}

	if qsize < wsize {
		qsize = wsize
	}

	pool := &Pool{
		queue:    make(chan Task, qsize),
		workers:  make(chan chan Task, wsize),
		shutdown: make(chan struct{}),
	}

	go pool.start()

	for i := 0; i < wsize; i++ {
		StartWorker(pool)
	}

	return pool
}

// Starts the scheduling.
func (p *Pool) start() {
	for {
		select {
		case worker := <-p.workers:
			task := <-p.queue
			worker <- task
		case <-p.shutdown:
			p.stopGracefully()
			return
		}
	}
}

func (p *Pool) isShutdown() error {
	select {
	case <-p.shutdown:
		return errSchedulerPoolStop
	default:
	}

	return nil
}

// Schedule push a task on queue.
func (p *Pool) Schedule(ctx context.Context, task Task) error {
	if err := p.isShutdown(); err != nil {
		return err
	}

	p.wg.Add(1)
	t := &taskWrapper{task, ctx}
	p.queue <- t
	return nil
}

// ScheduleWithTimeout try to push a task on queue, if timeout, return false.
func (p *Pool) ScheduleWithTimeout(ctx context.Context, timeout time.Duration, task Task) error {
	if err := p.isShutdown(); err != nil {
		return err
	}

	timer := time.NewTimer(timeout)
	t := &taskWrapper{task, ctx}

	select {
	case p.queue <- t:
		timer.Stop()
		return nil
	case <-timer.C:
		return ErrScheduleTimeout
	}
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.shutdown)
	})
}

func (p *Pool) stopGracefully() {
	// todo
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
