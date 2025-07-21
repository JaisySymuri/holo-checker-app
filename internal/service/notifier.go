package service

import (
	"fmt"
	"holo-checker-app/internal/controller"
	"holo-checker-app/internal/utility"
	"time"

	"github.com/sirupsen/logrus"
)


func Notify(videoInfos []utility.APIVideoInfo) error {
	var message string

	if len(videoInfos) == 0 {
		msg, err := makeNotFoundMessage()
		if err != nil {
			return err
		}
		message += msg + "\n"
	} else {
		for _, info := range videoInfos {
			msg, err := makeFoundMessage(info)
			if err != nil {
				return err
			}
			message += msg + "\n"
		}
	}

	// Send the message (to Telegram, WhatsApp, etc.)
	if err := controller.SendMessageToTelegram(utility.BotToken, utility.ChatID, message); err != nil {
		return err
	}
	if err := controller.SendMessageToWhatsApp(utility.PhoneNumber, utility.ApiKey, message); err != nil {
		return err
	}

	return nil
}

func makeFoundMessage(info utility.APIVideoInfo) (string, error) {
	var startTime time.Time
	var err error

	if info.StartScheduled != "" {
		startTime, err = time.Parse(time.RFC3339, info.StartScheduled)
		if err != nil {
			logrus.Debugf("Start Scheduled time for %s is not in RFC3339 format: %s", info.ID, info.StartScheduled)
			return "", fmt.Errorf("failed to parse StartScheduled time: %w", err)
		}
	} else {
		logrus.Debugf("Start Scheduled time for %s is empty, skipping parse", info.ID)
	}

	durationUntilStart := time.Until(startTime)

	message := fmt.Sprintf(
		"%s: Found '%s' with channel '%s'\nStarts/ed: %s\n",
		info.Status, info.TopicID, info.Channel.Name, FormatDuration(durationUntilStart),
	)

	// logrus.Debug("Debug: ", message)	

	return message, nil
}

func makeNotFoundMessage() (string, error) {
	message := "API: No 'Singing' stream scheduled."

	// logrus.Debug("Debug: ", message)

	return message, nil
}



func makeStartedMessage(info utility.APIVideoInfo) (string, error) {
    if info.ID == "" || info.Channel.Name == "" {
        return "", fmt.Errorf("missing video ID or channel name")
    }

    message := fmt.Sprintf(
        "%s is live! Watch now: https://www.youtube.com/watch?v=%s (channel: %s)",
        info.Title, info.ID, info.Channel.Name,
    )
    return message, nil
}

// ---------- Presentation layer ----------
type Notifier interface {
    Started(info utility.APIVideoInfo) error
}

// One implementation that uses your helper functions + logrus
type multiNotifier struct{}



func (multiNotifier) send(msg string) error {
    // Telegram
    if err := controller.SendMessageToTelegram(
        utility.BotToken, utility.ChatID, msg); err != nil {
        return fmt.Errorf("telegram: %w", err)
    }
    // WhatsApp
    if err := controller.SendMessageToWhatsApp(
        utility.PhoneNumber, utility.ApiKey, msg); err != nil {
        return fmt.Errorf("whatsapp: %w", err)
    }
    return nil
}

func (n multiNotifier) Started(info utility.APIVideoInfo) error {
    msg, err := makeStartedMessage(info)
    if err != nil {
        return err
    }
    return n.send(msg)
}

