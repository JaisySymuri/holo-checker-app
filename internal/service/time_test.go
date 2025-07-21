package service

// import (
// 	"holo-checker-app/internal/controller"
// 	"testing"
// 	"time"
// )



// func TestVideoTimeToStart(t *testing.T) {
// 	var fetcher controller.VideoFetcher = &controller.MockHolodexClient{} // note: must be a pointer

// 	println("TimeNow: ", TimeNow().Format("2006-01-02 15:04:05"))

// 	videos, err := fetcher.FetchVideos()
// 	if err != nil {
// 		t.Fatalf("FetchVideos() failed: %v", err)
// 	}

// 	for _, raw := range videos {
// 		video := &APIVideoInfo{raw} // Wrap utility.APIVideoInfo

// 		startTime, err := time.Parse(time.RFC3339, video.StartScheduled)
// 		if err != nil {
// 			t.Errorf("Failed to parse StartScheduled for video %s: %v", video.Title, err)
// 			continue
// 		}

// 		println("--------")
// 		println("Video Title       :", video.Title)
// 		println("StartScheduled    :", startTime.Format("2006-01-02 15:04:05"))
// 		println("Current TimeNow() :", TimeNow().Format("2006-01-02 15:04:05"))

// 		result := video.TimeToStart()
// 		t.Logf("Video: %s | TimeToStart(): %s", video.Title, result)

// 		if result == "" {
// 			t.Errorf("Expected non-empty TimeToStart() for video: %s", video.Title)
// 		}
// 	}
// }
