package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
	jobExecID, err := strconv.Atoi(os.Getenv("JOB_EXECUTION_ID"))
	if err != nil {
		log.Fatalln(err)
	}

	res, err := client.GetExecution(ctx, jobExecID)
	if err != nil {
		log.Fatalf("get job execution error: %+v", err)
	}
	log.Printf("getting job execution: id = %d", res.JobExecutionID)
}
