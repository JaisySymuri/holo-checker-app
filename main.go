// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/service"
	"holo-checker-app/internal/utility"
	"net"

	// "net/http"
	_ "net/http/pprof"
	"time"

	"github.com/getlantern/systray"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func checkSSH(ctx context.Context, host, port string) error {
	address := net.JoinHostPort(host, port)
	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func checkWithRetry(host, port string, attempts int, timeout time.Duration) bool {
	for i := 1; i <= attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		err := checkSSH(ctx, host, port)
		cancel()

		if err == nil {
			return true
		}

		logrus.Warnf("Attempt %d failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	return false
}

func main() {
	utility.SetLog()
	utility.SetEnv()

	host := "168.110.203.131"
	port := "22"

	if checkWithRetry(host, port, 3, 3*time.Second) {
		logrus.Warn("Linux server reachable. Exiting Windows app.")
		return
	}

	logrus.Error("Linux server unreachable. Holochecker Windows app running.")

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
