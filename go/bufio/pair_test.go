package bufio

import (
	"bytes"
	"strings"
	"testing"
)

var pairTests = []string{
	"",
	"{abc}",
	"{abc}sss",
	"{a[b]c}",
	"}{abc}",
}
var espectedpairTests = []string{
	"",
	"{abc}",
	"{abc}",
	"{a[b]c}",
	"{abc}",
}

func TestPairScanner(t *testing.T) {

	for n, test := range pairTests {
		buf := strings.NewReader(test)
		s := NewPairScanner(buf).SetDiscardLeading(true)
		p, err := s.ScanDelimiters("{}")
		if err != nil && n != 0 {
			t.Errorf("#%d: Scan error:%v\n", n, err)
			continue
		}

		if !bytes.Equal(p, []byte(espectedpairTests[n])) {
			t.Errorf("#%d: expected %q got %q", n, test, string(p))
		}
	}
}
