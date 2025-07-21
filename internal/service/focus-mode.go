package service

import (
	"fmt"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// FocusMode holds the ticker and a channel to signal stop.
type FocusMode struct {
	ticker   *time.Ticker
	stopChan chan struct{}
	poller   Poller
	notifier  Notifier   // NEW
}

type FetchByIDFn func(string) (*utility.APIVideoInfo, error)

// focusModes is a registry of active focus modes.
// It is protected by a mutex for concurrent access.
var (
	focusModes   = make(map[string]*FocusMode)
	focusModesMu sync.Mutex
)

func scheduleFocusMode(videos []utility.APIVideoInfo) {
	for _, video := range videos {
		if video.StartScheduled == "" {
			continue // Skip if empty
		}

		startTime, err := time.Parse(time.RFC3339, video.StartScheduled)
		if err != nil {
			logrus.Debugf("StartScheduled time for %s is not in RFC3339 format: %s", video.ID, video.StartScheduled)
			continue // Skip if not parseable
		}

		delay := time.Until(startTime)
		if delay <= 0 {
			continue // Skip if time has passed
		}

		go func(v utility.APIVideoInfo) {
			timer := time.NewTimer(delay)
			<-timer.C
			StartFocusMode(v, 2*time.Minute)
		}(video)
	}
}

/* ---------- Domain layer ---------- */

// Poller describes anything that can decide “keep polling or stop”.
type PollResult int

const (
    NotYet PollResult = iota
    Started
)

type Poller interface {
    Poll() (PollResult, *utility.APIVideoInfo, error)
}

// holodexPoller is one concrete strategy.
type holodexPoller struct {
	video utility.APIVideoInfo
	fetchByID  FetchByIDFn
}

func newHolodexPoller(video utility.APIVideoInfo, fetch FetchByIDFn) *holodexPoller {
    return &holodexPoller{
        video:     video,
        fetchByID: fetch,
    }
}

func (h *holodexPoller) Poll() (PollResult, *utility.APIVideoInfo, error) {
    v, err := h.fetchByID(h.video.ID)
    if err != nil {
        return NotYet, nil, err        // worker can log the error
    }
    if v.Status == "live" {
        return Started, v, nil
    }
    return NotYet, nil, nil
}

/* ---------- Worker ---------- */

func newFocusMode(interval time.Duration, p Poller, n Notifier) *FocusMode {
    return &FocusMode{
        ticker:   time.NewTicker(interval),
        stopChan: make(chan struct{}),
        poller:   p,
        notifier: n,
    }
}

func (fm *FocusMode) run(cleanup func()) {
    defer cleanup()
    defer fm.ticker.Stop()

    if fm.doPoll() { return }
    for {
        select {
        case <-fm.ticker.C:
            if fm.doPoll() { return }
        case <-fm.stopChan:
            logrus.Info("🛑 Focus mode stopped by caller")
            return
        }
    }
}

func (fm *FocusMode) doPoll() bool {
    res, info, err := fm.poller.Poll()
    if err != nil {
        logrus.Errorf("poll error: %v", err)
        return false
    }
    switch res {
    case Started:
        _ = fm.notifier.Started(*info)
        return true
    }
    return false
}

/* ---------- Scheduler ---------- */


// StartFocusMode registers and schedules a new focus‑mode job.
// interval is injected (e.g. 2*time.Minute in prod, 3*time.Second in tests).
func StartFocusMode(video utility.APIVideoInfo, interval time.Duration) {
    focusModesMu.Lock()
    defer focusModesMu.Unlock()
    if _, exists := focusModes[video.ID]; exists {
        fmt.Printf("Focus mode already running for %s\n", video.ID)
        return
    }

    p := newHolodexPoller(video, controller.RequestHolodexByID)
    n := multiNotifier{}                           
    fm := newFocusMode(interval, p, n)
    focusModes[video.ID] = fm

    go fm.run(func() {
        focusModesMu.Lock()
        delete(focusModes, video.ID)
        focusModesMu.Unlock()
    })
    logrus.Infof("🔎 Focus mode started for: %s [%s]", video.Title, video.ID)
}


func StopAllFocusModes() {
    focusModesMu.Lock()
    defer focusModesMu.Unlock()
    for link, fm := range focusModes {
        close(fm.stopChan)
        delete(focusModes, link)
        fmt.Printf("Focus mode stopped for %s\n", link)
    }
}



