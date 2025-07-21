package utility

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Retry attempts to execute the provided function up to a specified number of times,
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		logrus.Errorf("Attempt %d failed: %v. Retrying in %s...", i+1, err, sleep)
		time.Sleep(sleep)
	}
	return fmt.Errorf("all attempts failed")
}

func SetupTestEnv() {
	logrus.SetLevel(logrus.DebugLevel)
    SetEnv()    
}





