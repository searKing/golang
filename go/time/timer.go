package time

import (
	"time"

	"github.com/searKing/golang/go/sync/atomic"
)

// https://github.com/golang/go/issues/27169
// Timer to fix time: Timer.Stop documentation example easily leads to deadlocks
type Timer struct {
	*time.Timer
	chanConsumed atomic.Bool
}

// saw channel read, must be called after receiving value from timer chan
// an example case is AfterFunc bellow.
func (t *Timer) ChanConsumed() {
	t.chanConsumed.Store(true)
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
// Reset can be invoked anytime, which enhances std time.Reset
// that should be invoked only on stopped or expired timers with drained channels,
func (t *Timer) Reset(d time.Duration) bool {
	ret := t.Stop()
	if !ret && !t.chanConsumed.Load() {
		// drain the channel, prevents the Timer from blocking on Send to t.C by sendTime, t.C is reused.
		// The underlying Timer is not recovered by the garbage collector until the timer fires.
		// consume the channel only once for the channel can be triggered only one time at most before Stop is called.
		<-t.C
	}
	t.Timer.Reset(d)
	t.chanConsumed.Store(false)
	return ret
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		Timer: time.NewTimer(d),
	}
}

func WrapTimer(t *time.Timer) *Timer {
	return &Timer{
		Timer: t,
	}
}

func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

func AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{}
	t.Timer = time.AfterFunc(d, func() {
		t.ChanConsumed()
		f()
	})
	return t
}
