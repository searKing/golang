package sink

import (
	"github.com/searKing/golang/go/util/function/consumer"
)

type TODO struct {
	consumer.TODO
}

func (_ *TODO) Begin(size int) {
	return
}

func (_ *TODO) End() {
	return
}

func (_ *TODO) CancellationRequested() bool {
	return false
}
