package gamma

import "testing"

func TestNewCancelReader(t *testing.T) {
	_, _ = NewCancelReader(nil)
}
