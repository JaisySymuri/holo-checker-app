package service

import (
	"fmt"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// FocusMode holds the ticker and a channel to signal stop.
type FocusMode struct {
	ticker   *time.Ticker
	stopChan chan struct{}
	poller   Poller
	notifier Notifier // NEW
}

type FetchByIDFn func(string) (*utility.APIVideoInfo, error)

// focusModes is a registry of active focus modes.
// It is protected by a mutex for concurrent access.
var (
	focusModes   = make(map[string]*FocusMode)
	focusModesMu sync.Mutex
)

func (km *KaraokeManager) AddScheduledVideo(v utility.APIVideoInfo) {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.scheduledVideos[v.ID] = v
}

func (km *KaraokeManager) GetScheduledVideos() []utility.APIVideoInfo {
	km.mu.RLock()
	defer km.mu.RUnlock()
	videos := make([]utility.APIVideoInfo, 0, len(km.scheduledVideos))
	for _, v := range km.scheduledVideos {
		videos = append(videos, v)
	}
	return videos
}

func (km *KaraokeManager) RemoveScheduledVideo(id string) {
	km.mu.Lock()
	defer km.mu.Unlock()
	delete(km.scheduledVideos, id)
}

func scheduleFocusMode(km *KaraokeManager, videos []utility.APIVideoInfo) {
	// First, add ALL valid videos to the scheduled list
	for _, video := range videos {
		if video.StartScheduled == "" {
			continue
		}

		startTime, err := time.Parse(time.RFC3339, video.StartScheduled)
		if err != nil {
			logrus.Debugf("StartScheduled time for %s is not in RFC3339 format: %s", video.ID, video.StartScheduled)
			continue
		}

		km.AddScheduledVideo(video)
		logrus.Infof("Video %s scheduled to start focus mode at %s", video.Channel.Name, startTime.Format(time.RFC3339))
	}

	// Then, set timers for each scheduled video
	for _, video := range videos {	
		startTime, err := time.Parse(time.RFC3339, video.StartScheduled)
		if err != nil {
			continue
		}

		delay := time.Until(startTime)

		go func(v utility.APIVideoInfo) {
			timer := time.NewTimer(delay)
			<-timer.C
			StartFocusMode(v, 2*time.Minute)			
		}(video)
	}

	scheduled := km.GetScheduledVideos()

	count := len(scheduled)
	names := make([]string, 0, count)
	for _, v := range scheduled {
		names = append(names, v.Channel.Name)
	}

	logrus.Infof("scheduleFocusMode: %d scheduled videos", count)

	if count > 0 {
		logrus.Infof("scheduleFocusMode: Scheduled channels: %s", strings.Join(names, ", "))
	} else {
		logrus.Info("scheduleFocusMode: No scheduled videos")
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
	video     utility.APIVideoInfo
	fetchByID FetchByIDFn
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
		return NotYet, nil, err // worker can log the error
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

func (fm *FocusMode) run() {
	defer fm.ticker.Stop()

	if fm.doPoll() {
		return
	}

	for {
		select {
		case <-fm.ticker.C:
			if fm.doPoll() {
				return
			}
		case <-fm.stopChan:
			logrus.Info("🛑 Focus mode stopped by caller")
			return
		}
	}
}

func (fm *FocusMode) Stop(videoID string) {
	close(fm.stopChan)
	logrus.Infof("🛑 Focus mode stop requested for: %s", videoID)

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

	go fm.run()
	logrus.Infof("🔎 Focus mode started for: %s [%s]", video.Title, video.ID)
}

func StopAllFocusModes() {
	focusModesMu.Lock()
	defer focusModesMu.Unlock()

	for id, fm := range focusModes {
		fm.Stop(id)
		delete(focusModes, id)
	}
}
