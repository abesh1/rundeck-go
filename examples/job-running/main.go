package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/a6e5h1/rundeck-go"
)

var (
	client *rundeck.Client
)

func init() {
	apiVersion, err := strconv.Atoi(os.Getenv("API_VERSION"))
	if err != nil {
		log.Fatalln(err)
	}

	client = rundeck.NewClient(&rundeck.NewClientInput{
		BasePath:   os.Getenv("BASE_PATH"),
		Token:      os.Getenv("TOKEN"),
		APIVersion: apiVersion,
		Project:    os.Getenv("PROJECT"),
	})
}

func main() {
	ctx := context.Background()

	input := &rundeck.RunJobInput{
		JobID:   os.Getenv("JOB_ID"),
		Options: nil,
	}

	if unixRunAt := os.Getenv("UNIX_RUN_AT"); unixRunAt != "" {
		unixRunAt2, err := strconv.Atoi(unixRunAt)
		if err != nil {
			log.Fatalf("cast unix run at error: %+v", err)
		}
		runAt := time.Unix(int64(unixRunAt2), 0)
		input.RunAt = &runAt
	}

	res, err := client.RunJob(ctx, input)
	if err != nil {
		log.Fatalf("run job error: %+v", err)
	}
	log.Printf("running job execution: id = %d", res.JobExecutionID)
}
