package controller

import (
	"encoding/json"
	"fmt"
	"holo-checker-app/internal/utility"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type HolodexAPIClient struct {
	BaseURL string
	xApiKey string
	Client  *http.Client
}

// NewAPIClient constructs a new Holodex API client.
func NewAPIClient(apiKey string) *HolodexAPIClient {
	return &HolodexAPIClient{
		BaseURL: "https://holodex.net/api/v2/live",
		xApiKey: utility.XApiKey,
		Client:  &http.Client{},
	}
}

func (c *HolodexAPIClient) FetchVideos() ([]utility.APIVideoInfo, error) {
	types := []string{"stream", "placeholder"}
	topics := []string{"singing", "Marshmallow"}
	var allVideos []utility.APIVideoInfo

	for _, topic := range topics {
		for _, videoType := range types {
			videos, err := c.fetchVideosByTopicAndType(topic, videoType)
			if err != nil {
				return nil, err
			}
			allVideos = append(allVideos, videos...)
		}
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugf("FetchVideos: Fetched %d videos", len(allVideos))
		for _, v := range allVideos {
			fmt.Printf("%s🎵 %s (%s) by %s\n", v.TopicID, v.Title, v.Status, v.Channel.Name)
			fmt.Printf("   YouTube: https://www.youtube.com/watch?v=%s\n", v.ID)
			fmt.Println("--------------------------------------------------")
		}
	}

	return allVideos, nil
}

// Helper: fetch videos for one (topic, type)
func (c *HolodexAPIClient) fetchVideosByTopicAndType(topic, videoType string) ([]utility.APIVideoInfo, error) {
	params := url.Values{}
	params.Set("org", "Hololive")
	params.Set("topic", topic)
	params.Set("status", "new,upcoming,live")
	params.Set("type", videoType)
	params.Set("limit", "50")

	fullURL := fmt.Sprintf("%s?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-APIKEY", c.xApiKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var videos []utility.APIVideoInfo
	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return nil, err
	}

	return videos, nil
}

func RequestHolodexByID(videoID string) (*utility.APIVideoInfo, error) {
	// Prepare request
	baseURL := "https://holodex.net/api/v2/live"
	params := url.Values{}
	params.Set("id", videoID)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("X-APIKEY", utility.XApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Decode result
	var videos []utility.APIVideoInfo
	if err := json.NewDecoder(resp.Body).Decode(&videos); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no video found for ID: %s", videoID)
	}

	video := videos[0]
	logrus.Infof("🎯 Focus check: %s [%s] - Status: %s", video.Title, video.ID, video.Status)

	return &video, nil
}
