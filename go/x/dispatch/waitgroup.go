package dispatch

type WaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}

var nullWG = &emptyWG{}

type emptyWG struct {
	WaitGroup
}

func (wg *emptyWG) Add(delta int) {
	return
}

// Done decrements the waitGroup counter by one.
func (wg *emptyWG) Done() {
	return
}

// Wait blocks until the waitGroup counter is zero.
func (wg *emptyWG) Wait() {
	return
}
