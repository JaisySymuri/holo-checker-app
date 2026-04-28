//go:build linux

package service

import (
	"context"
	"holo-checker-app/internal/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func Run(km *KaraokeManager, fetcher controller.VideoFetcher) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartMonitoring(ctx, km, fetcher)
	logrus.Info("Running in headless Linux mode")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	sig := <-signalCh
	logrus.Infof("Received signal %s, shutting down", sig)
	SetRunning(false)
	StopAllFocusModes()
	cancel()
}
