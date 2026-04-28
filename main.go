// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"github.com/sirupsen/logrus"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/service"
	"holo-checker-app/internal/utility"
	_ "net/http/pprof"
)

func main() {
	utility.SetLog()
	utility.SetEnv()

	km := service.NewKaraokeManager()
	apiClient := controller.NewAPIClient(utility.XApiKey)

	logrus.Info("checkHolodex started. Connecting to internet...")
	service.Run(km, apiClient)
}
