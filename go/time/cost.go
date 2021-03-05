package time

import (
	"fmt"
	"io"
	"sort"
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
	costs    []time.Duration
	// Lesser reports whether the element with duration i
	// must sort before the element with duration j.
	// sorting in decreasing order of cost if Lesser is nil
	// behaves like Less in sort.Interface
	Lesser func(i time.Duration, j time.Duration) bool
}

func (c *CostTick) Reset() {
	c.points = nil
	c.messages = nil
	c.costs = nil
}

func (c *CostTick) Tick(msg string) {
	if msg == "" {
		caller, file, line := runtime_.GetShortCallerFuncFileLine(1)
		msg = fmt.Sprintf("%s() %s:%d", caller, file, line)
	}
	c.points = append(c.points, time.Now())
	c.messages = append(c.messages, msg)
	if len(c.costs) == 0 || len(c.points) == 1 {
		c.costs = append(c.costs, 0)
	} else {
		c.costs = append(c.costs, c.points[len(c.points)-1].Sub(c.points[len(c.points)-2]))
	}
}

func (c CostTick) String() string {
	var buf strings.Builder
	var scanned bool
	c.Walk(func(idx int, msg string, cost time.Duration, at time.Time) (next bool) {
		if scanned {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("#%d, msg: %s, cost %s, at %s", idx, msg, cost, at))
		scanned = true
		return true
	})
	return buf.String()
}

func (c CostTick) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, c.String())
		}
		fallthrough
	case 's', 'q':
		var costs []string
		c.Walk(func(idx int, msg string, cost time.Duration, at time.Time) (next bool) {
			costs = append(costs, fmt.Sprintf("%s(%s)", msg, cost))
			return true
		})
		_, _ = fmt.Fprintf(s, "%v", costs)
	}
}

func (c CostTick) Costs() []time.Duration {
	var costs []time.Duration
	c.Walk(func(idx int, msg string, cost time.Duration, at time.Time) (next bool) {
		costs = append(costs, cost)
		return true
	})
	return costs
}

// Walk iterate costs
// Walk stop if f(...) returns false.
func (c CostTick) Walk(f func(idx int, msg string, cost time.Duration, at time.Time) (next bool)) {
	if f == nil {
		return
	}

	for i, p := range c.points {
		if f(i, c.messages[i], c.costs[i], p) {
			continue
		}
		return
	}
}

// sorting in decreasing order of cost.
func (c CostTick) Len() int { return len(c.points) }

func (c CostTick) Less(i, j int) bool {
	if c.Lesser != nil {
		return c.Lesser(c.costs[i], c.costs[j])
	}
	return c.costs[i] > c.costs[j]
}

func (c *CostTick) Swap(i, j int) {
	c.points[i], c.points[j] = c.points[j], c.points[i]
	c.messages[i], c.messages[j] = c.messages[j], c.messages[i]
	c.costs[i], c.costs[j] = c.costs[j], c.costs[i]
}

// Sort is a convenience method: x.Sort() calls Sort(x).
func (c *CostTick) Sort() { sort.Sort(c) }
