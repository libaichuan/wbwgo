package timingwheel

import (
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	w := NewTimingWheel(time.Second, 2)

	for {
		select {
		case <-w.After(5 * time.Second):
			return
		}
	}
}
