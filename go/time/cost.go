package time

import (
	"fmt"
	"strings"
	"time"

	runtime_ "github.com/searKing/golang/go/runtime"
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

	flag int // properties
}

func (c *CostTick) Reset() {
	c.points = nil
	c.messages = nil
}

func (c *CostTick) Tick(msg string) {
	if msg == "" {
		caller, file, line := runtime_.GetShortCallerFuncFileLine(1)
		msg = fmt.Sprintf("%s() %s:%d", caller, file, line)
	}
	c.points = append(c.points, time.Now())
	c.messages = append(c.messages, msg)
}

func (c *CostTick) String() string {
	var buf strings.Builder
	var scanned bool
	c.Summary(func(idx int, msg string, cost time.Duration, at time.Time) {
		if scanned {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("#%d, msg: %s, cost %s, at %s", idx, msg, cost, at))
		scanned = true
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

	for i, p := range c.points {
		if i == 0 {
			f(i, c.messages[i], 0, p)
			continue
		}
		f(i, c.messages[i], p.Sub(c.points[i-1]), p)
	}
}
