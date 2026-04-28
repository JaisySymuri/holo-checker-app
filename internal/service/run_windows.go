//go:build windows

package service

import (
	"context"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"os"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

func Run(km *KaraokeManager, fetcher controller.VideoFetcher) {
	ctx, cancel := context.WithCancel(context.Background())
	StartMonitoring(ctx, km, fetcher)

	systray.Run(func() { onReady(km, fetcher, cancel) }, func() { onExit(cancel) })
}

func onReady(km *KaraokeManager, fetcher controller.VideoFetcher, cancel context.CancelFunc) {
	iconData, err := os.ReadFile("favicon.ico")
	if err != nil {
		logrus.Fatalf("Failed to read icon file: %v", err)
	}

	systray.SetIcon(iconData)
	systray.SetTitle("Holodex Checker")
	systray.SetTooltip("Holodex Checker")

	startMenuItem := systray.AddMenuItem("Start", "Start checking Holodex")
	pauseMenuItem := systray.AddMenuItem("Pause", "Pause checking Holodex")
	restartMenuItem := systray.AddMenuItem("Restart", "Restart checking Holodex")
	exitMenuItem := systray.AddMenuItem("Exit", "Exit the application")
	hideConsoleMenuItem := systray.AddMenuItem("Hide Console", "Hide the console window")
	stopFocusMode := systray.AddMenuItem("Stop focus", "Stopping focus mode for the earliest stream")

	apiClient := controller.NewAPIClient(utility.XApiKey)

	go func() {
		for {
			select {
			case <-startMenuItem.ClickedCh:
				if !IsRunning() {
					SetRunning(true)
					logrus.Info("checkHolodex started")
					TriggerMonitor(km, apiClient)
				}
			case <-pauseMenuItem.ClickedCh:
				if IsRunning() {
					SetRunning(false)
					logrus.Info("checkHolodex paused")
				}
			case <-restartMenuItem.ClickedCh:
				SetRunning(false)
				logrus.Info("checkHolodex restarting")
				time.Sleep(2 * time.Second)
				SetRunning(true)
				TriggerMonitor(km, apiClient)
			case <-hideConsoleMenuItem.ClickedCh:
				syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole").Call()
				logrus.Info("Console window hidden")
			case <-stopFocusMode.ClickedCh:
				StopAllFocusModes()
			case <-exitMenuItem.ClickedCh:
				cancel()
				SetRunning(false)
				StopAllFocusModes()
				systray.Quit()
				logrus.Info("Exiting...")
				return
			}
		}
	}()
}

func onExit(cancel context.CancelFunc) {
	cancel()
	SetRunning(false)
	logrus.Info("Application exited")
}
