// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"holo-checker-app/internal/service"
	"holo-checker-app/internal/utility"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/getlantern/systray"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	utility.SetLog()
	utility.SetEnv()

	
	km := service.NewKaraokeManager()

	logrus.Info("checkHolodex started. Connecting to internet...")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe("localhost:2112", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		service.Monitor(km)
	
		// Align to next 10-minute mark
		now := time.Now()
		first := now.Truncate(10 * time.Minute).Add(10 * time.Minute)
		logrus.Debugf("Next monitor run at %v", first)
		time.Sleep(time.Until(first))
	
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
	
		for service.Running {
			service.Monitor(km)
	
			next := time.Now().Add(10 * time.Minute).Truncate(10 * time.Minute)
			logrus.Debugf("Next monitor run at %v", next)
	
			<-ticker.C
		}
	}()
	

	systray.Run(func() { service.OnReady(km) }, service.OnExit)
}

