package time

import "time"

type Cost struct {
	start time.Time
}

func (c *Cost) Start() {
	c.start = time.Now()
}

func (c *Cost) Elapse() time.Duration {
	return time.Now().Sub(c.start)
}

func (c *Cost) End(f func(d time.Duration)) {
	if f != nil {
		f(c.Elapse())
	}

}
