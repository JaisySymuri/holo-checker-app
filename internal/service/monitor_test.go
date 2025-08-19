package service

import (
	"holo-checker-app/internal/mockdata"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	mockdata.GenerateHolodexJSON(10 * time.Second)

	// Prepare KaraokeManager
    km := &KaraokeManager{}

    // Use mock fetcher (loads from testdata/holodex.json)
    mockFetcher := &mockdata.MockFetcher{}

    // Call Monitor with the mock
    Monitor(km, mockFetcher)

    // Verify mock loaded videos
    if len(mockFetcher.Videos) == 0 {
        t.Fatal("expected mock videos to be loaded from testdata/holodex.json, got 0")
    }

    // Example assertion: ensure km has scheduled videos after Monitor
    scheduled := km.GetScheduledVideos()
    if len(scheduled) == 0 {
        t.Error("expected scheduled videos to be set in KaraokeManager, got 0")
    }

    // Optionally check that specific fields match expected values
    for _, v := range scheduled {
        if v.Title == "" {
            t.Errorf("scheduled video has empty title: %+v", v)
        }
    }
}
