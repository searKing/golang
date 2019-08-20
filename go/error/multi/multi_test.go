package multi_test

import (
	"fmt"
	"github.com/searKing/golang/go/error/multi"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		errs []error
		want error
	}{
		{nil, nil},
		{[]error{}, nil},
		{[]error{fmt.Errorf("foo")}, fmt.Errorf("foo")},
		{[]error{fmt.Errorf("foo"), fmt.Errorf("fun")}, multi.New(fmt.Errorf("foo"), fmt.Errorf("fun"))},
	}

	for _, tt := range tests {
		got := multi.New(tt.errs...)
		if got == tt.want {
			continue
		}
		if got == nil || tt.want == nil || got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}
