// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/service"
	"holo-checker-app/internal/utility"
	"os/exec"
	"runtime"

	// "net/http"
	_ "net/http/pprof"
	"time"

	"github.com/getlantern/systray"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func isHostReachable(host string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", host)
	}

	err := cmd.Run()
	return err == nil
}

func main() {
	utility.SetLog()
	utility.SetEnv()

	if isHostReachable("168.110.203.131") {
		logrus.Warn("Linux server reachable. Exiting Windows app.")
		return
	}

	km := service.NewKaraokeManager()
	apiClient := controller.NewAPIClient(utility.XApiKey)

	logrus.Info("checkHolodex started. Connecting to internet...")

	// go func() {
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	if err := http.ListenAndServe("localhost:2112", nil); err != nil {
	// 		panic(err)
	// 	}
	// }()

	go func() {
		service.Monitor(km, apiClient)

		// Align to next 10-minute mark
		now := time.Now()
		first := now.Truncate(10 * time.Minute).Add(10 * time.Minute)
		logrus.Debugf("Next monitor run at %v", first)
		time.Sleep(time.Until(first))

		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for service.Running {
			service.Monitor(km, apiClient)

			next := time.Now().Add(10 * time.Minute).Truncate(10 * time.Minute)
			logrus.Debugf("Next monitor run at %v", next)

			<-ticker.C
		}
	}()

	systray.Run(func() { service.OnReady(km) }, service.OnExit)
}
