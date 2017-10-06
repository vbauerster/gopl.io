package memo

import (
	"memotest"
	"testing"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := New(httpGetBody)
	defer m.Close()
	memotest.Concurrent(t, m)
}
