package consumer

import (
	"github.com/searKing/golang/go/util/object"
)

// Consumer represents an operation that accepts a single input argument and returns no
// result. Unlike most other functional interfaces, {@code Consumer} is expected
// to operate via side-effects.
type Consumer interface {
	Accept(value interface{})
}

var NopThenableConsumer = ThenableConsumer(func(value interface{}) {})

type ThenableConsumer func(value interface{})

func (c ThenableConsumer) Accept(value interface{}) {
	object.RequireNonNil(c)
	acceptFn := (func(value interface{}))(c)
	acceptFn(value)
}

func (c ThenableConsumer) AndThen(consumer Consumer) ThenableConsumer {
	return ThenableConsumer(
		func(value interface{}) {
			c.Accept(value)
			consumer.Accept(value)
		})
}

type EmptyConsumer interface {
	Run()
}
