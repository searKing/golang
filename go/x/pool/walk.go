package pool

import (
	"context"
	"sync"
)

type Walk struct {
	wg sync.WaitGroup

	walkError     error
	walkErrorOnce sync.Once
	walkErrorMu   sync.Mutex
}

// WalkFunc is the type of the function called for each task processed
// by Walk. The path argument contains the argument to Walk as a task.
//
// In the case of an error, the info argument will be nil. If an error
// is returned, processing stops.
type WalkFunc func(task interface{}) error

// Walk will consume all tasks parallel and block until ctx.Done() or taskChan is closed.
func (p *Walk) Walk(ctx context.Context, taskChan <-chan interface{}, procFn WalkFunc) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ctx, cancel := context.WithCancel(ctx)

		for {
			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}
				p.wg.Add(1)
				go func() {
					defer p.wg.Done()
					if err := procFn(task); err != nil {
						p.TrySetError(err)
						cancel()
						return
					}
				}()
			case <-ctx.Done():
				p.TrySetError(ctx.Err())
				return
			}
		}
	}()
}

// Wait blocks until the WaitGroup counter is zero.
func (p *Walk) Wait() {
	p.wg.Wait()
}

func (p *Walk) Error() error {
	p.walkErrorMu.Lock()
	defer p.walkErrorMu.Unlock()
	return p.walkError
}
func (p *Walk) TrySetError(err error) {
	p.walkErrorOnce.Do(func() {
		p.walkErrorMu.Lock()
		defer p.walkErrorMu.Unlock()
		p.walkError = err
	})
}
