package time

import (
	"fmt"
	"strings"
	"time"
)

type Cost struct {
	start time.Time
}

func (c *Cost) Start() {
	c.start = time.Now()
}

func (c *Cost) Elapse() time.Duration {
	return time.Now().Sub(c.start)
}

func (c *Cost) ElapseFunc(f func(d time.Duration)) {
	if f != nil {
		f(c.Elapse())
	}
}

type CostTick struct {
	points   []time.Time
	messages []string
}

func (c *CostTick) Reset() {
	c.points = nil
	c.messages = nil
}

func (c *CostTick) Tick(msg string) {
	c.points = append(c.points, time.Now())
	c.messages = append(c.messages, msg)
}

func (c *CostTick) String() string {
	var buf strings.Builder
	c.Summary(func(idx int, msg string, cost time.Duration, at time.Time) {
		buf.WriteString(fmt.Sprintf("#%d, msg: %s, cost %s, at %s", idx, msg, cost, at))
	})
	return buf.String()
}

func (c *CostTick) Summary(f func(idx int, msg string, cost time.Duration, at time.Time)) {
	if f == nil {
		return
	}
	if c == nil || len(c.points) == 0 {
		return
	}
	var start = c.points[0]

	for i, p := range c.points {
		f(i, c.messages[i], p.Sub(start), p)
	}
}
