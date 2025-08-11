package mockdata

import (
	"encoding/json"
	"fmt"
	"holo-checker-app/internal/utility"
	"os"
	"path/filepath"
)

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
