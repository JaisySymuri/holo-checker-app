package service

import (
	"context"
	"holo-checker-app/internal/controller"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var running atomic.Bool

func init() {
	running.Store(true)
}

func IsRunning() bool {
	return running.Load()
}

func SetRunning(value bool) {
	running.Store(value)
}

func StartMonitoring(ctx context.Context, km *KaraokeManager, fetcher controller.VideoFetcher) {
	go func() {
		runMonitorCycle(ctx, km, fetcher)
	}()
}

func TriggerMonitor(km *KaraokeManager, fetcher controller.VideoFetcher) {
	go Monitor(km, fetcher)
}

func runMonitorCycle(ctx context.Context, km *KaraokeManager, fetcher controller.VideoFetcher) {
	Monitor(km, fetcher)

	timer := time.NewTimer(time.Until(nextMonitorRun(time.Now())))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Monitoring loop stopped")
			return
		case <-timer.C:
			if IsRunning() {
				Monitor(km, fetcher)
			} else {
				logrus.Info("Monitor tick skipped because checking is paused")
			}

			next := nextMonitorRun(time.Now())
			logrus.Debugf("Next monitor run at %v", next)
			timer.Reset(time.Until(next))
		}
	}
}

func nextMonitorRun(now time.Time) time.Time {
	next := now.Truncate(10 * time.Minute).Add(10 * time.Minute)
	logrus.Debugf("Next monitor run at %v", next)
	return next
}
