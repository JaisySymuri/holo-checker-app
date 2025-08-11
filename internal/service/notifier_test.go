package service

import (
	"fmt"
	"holo-checker-app/internal/mockdata"

	"testing"
)

func TestMakeFoundMessage(t *testing.T) {
	videos, err := mockdata.LoadMockHolodexData()
	if err != nil {
		t.Fatalf("failed to load mock data: %v", err)
	}

	for i, video := range videos {
		msg, err := makeFoundMessage(video)
		if err != nil {
			t.Errorf("error on video %d (%s): %v", i, video.ID, err)
			continue
		}

		fmt.Println("────────────────────────────────────────")
		fmt.Printf("🎬 Video %d: %s\n", i+1, video.Title)
		fmt.Println(msg)
	}
}
