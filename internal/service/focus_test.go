package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"holo-checker-app/internal/utility"
)

// Mock Poller
type mockPoller struct{}

func (m mockPoller) Poll() (PollResult, *utility.APIVideoInfo, error) {
	return NotYet, nil, nil
}

// Mock Notifier
type mockNotifier struct{}

func (m mockNotifier) Started(info utility.APIVideoInfo) error {
	return nil
}

func TestNewFocusMode(t *testing.T) {
	interval := 2 * time.Second
	p := mockPoller{}
	n := mockNotifier{}

	fm := newFocusMode(interval, p, n)

	assert.NotNil(t, fm, "FocusMode should not be nil")
	assert.NotNil(t, fm.ticker, "Ticker should be initialized")
	assert.NotNil(t, fm.stopChan, "Stop channel should be initialized")

	// Ensure the ticker has the correct duration
	// NOTE: time.Ticker doesn't expose the interval directly,
	// but we can check if it's not nil (since interval is used in its creation).
	assert.Implements(t, (*Poller)(nil), fm.poller)
	assert.Implements(t, (*Notifier)(nil), fm.notifier)
}

func TestNewFocusMode_PrintTicker(t *testing.T) {
	interval := 2 * time.Second
	p := mockPoller{}
	n := mockNotifier{}

	fm := newFocusMode(interval, p, n)

	assert.NotNil(t, fm)
	assert.NotNil(t, fm.ticker)
	assert.NotNil(t, fm.stopChan)

	// Simulate ticker for 15 seconds
	stopAfter := 15 * time.Second
	stopTimer := time.NewTimer(stopAfter)

	t.Logf("Starting ticker test: printing every %v, stopping after %v", interval, stopAfter)

	go func() {
		for {
			select {
			case <-fm.ticker.C:
				t.Log("Ticker fired: doing work...")
			case <-stopTimer.C:
				t.Log("Stopping ticker after 15 seconds.")
				fm.ticker.Stop()
				close(fm.stopChan)
				return
			}
		}
	}()

	// Wait until stop timer triggers
	<-stopTimer.C
}
