package gamma

import (
	"sync/atomic"
	"time"
)

// TimerEvent is sent to the terminal's event channel when a timer fires.
type TimerEvent struct {
	// ID identifies the timer that fired.
	ID int
}

// Timer represents a timer that sends [TimerEvent]s to a [Terminal].
type Timer struct {
	id   int
	stop chan struct{}
}

// Stop cancels the timer. It is safe to call multiple times.
func (t *Timer) Stop() {
	select {
	case t.stop <- struct{}{}:
	default:
	}
}

// ID returns the timer's identifier.
func (t *Timer) ID() int {
	return t.id
}

var nextTimerID atomic.Int64

// After creates a timer that fires a single [TimerEvent] after the given
// duration. The returned [Timer] can be stopped before it fires.
func (t *Terminal) After(d time.Duration) *Timer {
	id := int(nextTimerID.Add(1))
	timer := &Timer{id: id, stop: make(chan struct{}, 1)}
	go func() {
		tm := time.NewTimer(d)
		defer tm.Stop()
		select {
		case <-tm.C:
			t.SendEvent(TimerEvent{ID: id})
		case <-timer.stop:
		case <-t.donec:
		}
	}()
	return timer
}

// Every creates a timer that fires a [TimerEvent] repeatedly at the given
// interval. The returned [Timer] must be stopped when no longer needed.
func (t *Terminal) Every(d time.Duration) *Timer {
	id := int(nextTimerID.Add(1))
	timer := &Timer{id: id, stop: make(chan struct{}, 1)}
	go func() {
		tk := time.NewTicker(d)
		defer tk.Stop()
		for {
			select {
			case <-tk.C:
				t.SendEvent(TimerEvent{ID: id})
			case <-timer.stop:
				return
			case <-t.donec:
				return
			}
		}
	}()
	return timer
}
