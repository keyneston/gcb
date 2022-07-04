package gcb

import "time"

// cleanTimer attempts to correctly drain the timer and prepare it for GC.
func cleanTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}
