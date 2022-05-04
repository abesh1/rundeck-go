package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/a6e5h1/rundeck-go"

	"github.com/google/uuid"
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

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("new uuid error: %+v", err)
	}

	res, err := client.ImportJob(ctx, &rundeck.ImportJobInput{
		Jobs: []*rundeck.JobInput{
			{
				UUID:               id.String(),
				Group:              "/sample/1",
				Name:               "sample job",
				Description:        "this is a sample job",
				LogLevel:           rundeck.LogLevelInfo,
				ExecutionEnabled:   true,
				MultipleExecutions: true,
				Schedule: &rundeck.ScheduleInput{
					Month:   "*",
					Hour:    12,
					Minute:  0,
					Second:  0,
					Weekday: "*",
					Year:    "*",
				},
				Sequence: &rundeck.JobSequenceInput{
					Commands: []string{
						"echo Hello!",
					},
					KeepGoing: false,
					Strategy:  rundeck.JobSequenceStrategyNodeFirst,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("import job error: %+v", err)
	}

	if len(res.Succeeded) > 0 {
		for _, j := range res.Succeeded {
			log.Printf("import succeeded job: id = %s, name = %s, permalink = %s", j.ID, j.Name, j.Permalink)
		}
	}
	if len(res.Failed) > 0 {
		for _, j := range res.Failed {
			log.Printf("import failed job: id = %s, name = %s, permalink = %s", j.ID, j.Name, j.Permalink)
		}
	}
	if len(res.Skipped) > 0 {
		for _, j := range res.Skipped {
			log.Printf("import skipped job: id = %s, name = %s, permalink = %s", j.ID, j.Name, j.Permalink)
		}
	}
}
