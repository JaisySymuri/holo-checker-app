package controller

import (
	"context"
	"fmt"
	"holo-checker-app/internal/utility"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskScheduler struct {
	client      *cloudtasks.Client
	projectID   string
	locationID  string
	queueID     string
	endpointURL string // your Cloud Run endpoint
}

func NewTaskScheduler(ctx context.Context, projectID, locationID, queueID, endpointURL string) (*TaskScheduler, error) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &TaskScheduler{
		client:      client,
		projectID:   projectID,
		locationID:  locationID,
		queueID:     queueID,
		endpointURL: endpointURL,
	}, nil
}

func (ts *TaskScheduler) ScheduleVideoTask(ctx context.Context, videoID string, runAt time.Time) error {
	// Task will POST to Cloud Run
	req := &taskspb.CreateTaskRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", ts.projectID, ts.locationID, ts.queueID),
		Task: &taskspb.Task{
			ScheduleTime: timestamppb.New(runAt),
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        ts.endpointURL,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: []byte(fmt.Sprintf(`{"video_id":"%s"}`, videoID)),
				},
			},
		},
	}

	_, err := ts.client.CreateTask(ctx, req)
	return err
}

func ScheduleAllVideos(videos []utility.APIVideoInfo) error {
	ctx := context.Background()
	scheduler, err := NewTaskScheduler(ctx, "your-project-id", "us-central1", "my-task-queue", "https://your-cloud-run-url.com/handle")
	if err != nil {
		return err
	}

	for _, video := range videos {
		if video.StartScheduled == "" {
			continue
		}

		scheduledTime, err := time.Parse(time.RFC3339, video.StartScheduled)
		if err != nil {
			continue
		}

		err = scheduler.ScheduleVideoTask(ctx, video.ID, scheduledTime)
		if err != nil {
			fmt.Printf("❌ Failed to schedule task for %s: %v\n", video.Title, err)
		} else {
			fmt.Printf("✅ Scheduled task for %s at %s\n", video.Title, scheduledTime)
		}
	}

	return nil
}
