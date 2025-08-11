package service

import (
	"holo-checker-app/internal/utility"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utility.SetupTestEnv()
	os.Exit(m.Run())
}
