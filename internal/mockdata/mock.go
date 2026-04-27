package mockdata

import (
	"encoding/json"
	"fmt"
	"holo-checker-app/internal/utility"
	"os"
	"path/filepath"
	"time"
)

type MockFetcher struct {
	Videos []utility.APIVideoInfo
	Err    error
}

func (m *MockFetcher) FetchVideos() ([]utility.APIVideoInfo, error) {
	m.Videos, m.Err = LoadMockHolodexData()
	return m.Videos, m.Err
}

type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Org         string `json:"org"`
	Suborg      string `json:"suborg"`
	Type        string `json:"type"`
	Photo       string `json:"photo"`
	EnglishName string `json:"english_name"`
}

type Video struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Type           string  `json:"type"`
	TopicID        string  `json:"topic_id"`
	PublishedAt    string  `json:"published_at"`
	AvailableAt    string  `json:"available_at"`
	Duration       int     `json:"duration"`
	Status         string  `json:"status"`
	StartScheduled string  `json:"start_scheduled"`
	LiveViewers    int     `json:"live_viewers"`
	Channel        Channel `json:"channel"`
}

func GenerateHolodexJSON(delay time.Duration) {
	// Read JSON file
	root, err := findProjectRoot()
	if err != nil {
		fmt.Println("Error finding project root:", err)
		return
	}
	filePath := filepath.Join(root, "testdata", "holodex.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse JSON
	var videos []Video
	if err := json.Unmarshal(data, &videos); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Get Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading Jakarta timezone:", err)
		return
	}

	// Calculate new time: now + delay (Jakarta local), then convert to UTC
	newTimeJakarta := time.Now().In(loc).Add(delay)
	newTimeUTC := newTimeJakarta.UTC()

	// Update all videos' start_scheduled
	for i := range videos {
		videos[i].StartScheduled = newTimeUTC.Format("2006-01-02T15:04:05.000Z")
	}

	// Save JSON back
	modified, err := json.MarshalIndent(videos, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	if err := os.WriteFile(filePath, modified, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated start_scheduled to:", newTimeUTC.Format(time.RFC3339Nano))
}

// LoadMockHolodexData always looks for the file in <project_root>/testdata/holodex.json
func LoadMockHolodexData() ([]utility.APIVideoInfo, error) {
	var videos []utility.APIVideoInfo

	// Use absolute path relative to project root
	root, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("cannot determine project root: %w", err)
	}
	mockPath := filepath.Join(root, "testdata", "holodex.json")

	fmt.Println("📁 Looking for testdata at:", mockPath)

	data, err := os.ReadFile(mockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read holodex.json: %w", err)
	}

	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return videos, nil
}

// LoadMockHolodexDataByID returns the first APIVideoInfo in testdata/holodex.json
// whose ID matches videoID.  If no match is found, it returns (nil, os.ErrNotExist).
func LoadMockHolodexDataByID(videoID string) (*utility.APIVideoInfo, error) {
	// Re‑use the existing loader so we keep one source of truth.
	videos, err := LoadMockHolodexData()
	if err != nil {
		return nil, err
	}

	for i := range videos {
		if videos[i].ID == videoID { // or videos[i].VideoID / YoutubeID – adjust field name
			return &videos[i], nil
		}
	}
	return nil, fmt.Errorf("video %q not found in mock data: %w", videoID, os.ErrNotExist)
}

// findProjectRoot walks up to find go.mod to locate root directory
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found in any parent directory")
}
