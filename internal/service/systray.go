package service

import (
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"os"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

var Running bool = true

func OnReady(km *KaraokeManager) {
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
				if !Running {
					Running = true
					logrus.Info("checkHolodex started")
					go Monitor(km, apiClient)
				}
			case <-pauseMenuItem.ClickedCh:
				if Running {
					Running = false
					logrus.Info("checkHolodex paused")
				}
			case <-restartMenuItem.ClickedCh:
				Running = false
				logrus.Info("checkHolodex restarting")
				time.Sleep(2 * time.Second)
				Running = true
				go Monitor(km, apiClient)
			case <-hideConsoleMenuItem.ClickedCh:
				syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole").Call()
				logrus.Info("Console window hidden")
			case <-stopFocusMode.ClickedCh:
				StopAllFocusModes()
			case <-exitMenuItem.ClickedCh:
				systray.Quit()
				logrus.Info("Exiting...")
				return
			}
		}
	}()
}

func OnExit() {
	logrus.Info("Application exited")
}
