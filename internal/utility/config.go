package utility

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func SetEnv() {
	// First try local path
	localPath, _ := filepath.Abs(".env")
	logrus.Debugf("Trying to load .env from: %s", localPath)

	err := godotenv.Load(".env")
	if err != nil {
		// Fallback to project root (up 2 dirs from internal/utility)
		fallbackPath := filepath.Join("..", "..", ".env")
		absFallback, _ := filepath.Abs(fallbackPath)
		logrus.Debugf("Fallback: trying to load .env from: %s", absFallback)

		err = godotenv.Load(fallbackPath)
		if err != nil {
			logrus.Fatalf("Failed to load .env from both paths. Last error: %v", err)
		}
	}

	// Set environment variables
	BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	ChatID = os.Getenv("TELEGRAM_CHAT_ID")
	PhoneNumber = os.Getenv("WHATSAPP_PHONE_NUMBER")
	ApiKey = os.Getenv("WHATSAPP_API_KEY")
	XApiKey = os.Getenv("XAPIKEY")
}


// Custom Log Formatter
type SimpleFormatter struct{}

func (f *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timeFormat := entry.Time.Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("%s %s\n", timeFormat, entry.Message)
	return []byte(message), nil
}

// Detect if a console is attached
func consoleAttached() bool {
    // GetStdHandle returns INVALID_HANDLE_VALUE or 0 if no console is attached
    h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
    return err == nil && h != 0 && h != syscall.InvalidHandle
}

func isGoRun() bool {
    exePath, err := os.Executable()
    if err != nil {
        return false // default to compiled mode
    }
    return strings.Contains(exePath, os.TempDir()) // "go run" builds in temp dir
}

func SetLog() {
    logrus.SetFormatter(&SimpleFormatter{})

    logFile, err := os.OpenFile("debug2.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatalf("Failed to open log file: %v", err)
    }

    if consoleAttached() {
        multiWriter := io.MultiWriter(os.Stdout, logFile)
        logrus.SetOutput(multiWriter)
        logrus.Debug("Console detected: logging to console + file")
    } else {
        logrus.SetOutput(logFile)
        logrus.Debug("No console detected: logging to file only")
    }

    // Set log level based on execution mode
    if isGoRun() {
        logrus.SetLevel(logrus.DebugLevel)
        logrus.Debug("Running via go run: DebugLevel enabled")
    } else {
        logrus.SetLevel(logrus.InfoLevel)
        logrus.Info("Running compiled binary: InfoLevel enabled")
    }

    logrus.Info("Logger initialized")
}


