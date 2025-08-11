package service

import (
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var appStartTime = time.Now()

func AppStartTime() time.Time {
	return appStartTime
}

type ChangeChecker interface {
	ShouldNotify(oldStreams, newStreams []utility.APIVideoInfo) bool
}

type DefaultChangeChecker struct{}

func (c DefaultChangeChecker) ShouldNotify(oldStreams, newStreams []utility.APIVideoInfo) bool {
	return hasNewIDs(oldStreams, newStreams) || isWithinFirst5Minutes()
}

// Check if there is a new ID
func hasNewIDs(oldStreams, newStreams []utility.APIVideoInfo) bool {
	oldIDs := make(map[string]struct{})
	for _, s := range oldStreams {
		oldIDs[s.ID] = struct{}{}
	}

	for _, s := range newStreams {
		if _, exists := oldIDs[s.ID]; !exists {
			return true
		}
	}
	return false
}

// Check if current time is within the first 5 minutes of the hour
func isWithinFirst5Minutes() bool {
	now := time.Now()
	return now.Minute() >= 0 && now.Minute() < 5
}

type KaraokeManager struct {
	streams         []utility.APIVideoInfo
	scheduledVideos map[string]utility.APIVideoInfo // key by ID or something unique
	mu              sync.RWMutex
}

func Monitor(km *KaraokeManager) {
	apiClient := controller.NewAPIClient(utility.XApiKey)
	checker := DefaultChangeChecker{}
	var newStreams []utility.APIVideoInfo

	err := utility.Retry(30, 10*time.Second, func() error {
		var err error
		newStreams, err = apiClient.FetchVideos()
		return err
	})
	if err != nil {
		logrus.Error("FetchVideos failed after retries: ", err)
		return
	}

	// Filter only Hololive streams
	hololiveStreams := FilterStreams(newStreams, IsHololive)

	handleStreamUpdate(km, checker, hololiveStreams)

	scheduled := km.GetScheduledVideos()

	count := len(scheduled)
	names := make([]string, 0, count)
	for _, v := range scheduled {
		names = append(names, v.Channel.Name)
	}

	logrus.Infof("Monitor: %d scheduled videos", count)

	if count > 0 {
		logrus.Infof("Monitor: Scheduled channels: %s", strings.Join(names, ", "))
	} else {
		logrus.Info("Monitor: No scheduled videos")
	}
}

func handleStreamUpdate(km *KaraokeManager, checker ChangeChecker, newStreams []utility.APIVideoInfo) {
	oldStreams := km.GetStreams()

	// Explicit first run condition
	if time.Since(AppStartTime()) < time.Minute {
		logrus.Info("First run detected, calling Notify and scheduling FocusMode (even if empty).")
		Notify(newStreams)
		km.SetStreams(newStreams)
		scheduleFocusMode(km, newStreams)
		return
	}

	// Usual logic for subsequent runs
	if checker.ShouldNotify(oldStreams, newStreams) {
		logrus.Info("Condition met, calling Notify and scheduling FocusMode...")
		Notify(newStreams)
		km.SetStreams(newStreams)
		scheduleFocusMode(km, newStreams)
	} else {
		km.SetStreams(newStreams)
		logrus.Info("No new streams and outside forced window, skipping Notify.")
	}
}

func NewKaraokeManager() *KaraokeManager {
	return &KaraokeManager{
		streams:         make([]utility.APIVideoInfo, 0),
		scheduledVideos: make(map[string]utility.APIVideoInfo),
	}
}

func (km *KaraokeManager) SetStreams(newStreams []utility.APIVideoInfo) {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.streams = newStreams
}

func (km *KaraokeManager) GetStreams() []utility.APIVideoInfo {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return append([]utility.APIVideoInfo{}, km.streams...)
}

func FilterStreams(streams []utility.APIVideoInfo, predicate func(utility.APIVideoInfo) bool) []utility.APIVideoInfo {
	var filtered []utility.APIVideoInfo
	for _, stream := range streams {
		if predicate(stream) {
			filtered = append(filtered, stream)
		}
	}
	return filtered
}

func IsHololive(stream utility.APIVideoInfo) bool {
	return stream.Channel.Org == "Hololive"
}
