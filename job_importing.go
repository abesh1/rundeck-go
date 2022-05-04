package rundeck

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
)

type JobDefinition struct {
	DefaultTab         string       `yaml:"defaultTab"`
	Description        string       `yaml:"description"`
	ExecutionEnabled   bool         `yaml:"executionEnabled"`
	Group              string       `yaml:"group"`
	ID                 string       `yaml:"id"`
	UUID               string       `yaml:"uuid"`
	LogLevel           string       `yaml:"loglevel"`
	MultipleExecutions bool         `yaml:"multipleExecutions"`
	Name               string       `yaml:"name"`
	NodeFilterEditable bool         `yaml:"nodeFilterEditable"`
	Schedule           *JobSchedule `yaml:"schedule"`
	ScheduleEnabled    bool         `yaml:"scheduleEnabled"`
	Sequence           *JobSequence `yaml:"sequence"`
}

type JobSchedule struct {
	Month   string              `yaml:"month"`
	Time    *JobScheduleTime    `yaml:"time"`
	Weekday *JobScheduleWeekday `yaml:"weekday"`
	Year    string              `yaml:"year"`
}

type JobScheduleTime struct {
	Hour    string `yaml:"hour"`
	Minute  string `yaml:"minute"`
	Seconds string `yaml:"seconds"`
}

type JobScheduleWeekday struct {
	Day string `yaml:"day"`
}

type JobSequence struct {
	Commands  []*JobCommand `yaml:"commands"`
	KeepGoing bool          `yaml:"keepgoing"`
	Strategy  string        `yaml:"strategy"`
}

type JobCommand struct {
	Exec string `yaml:"exec"`
}

const (
	LogLevelDebug   = "DEBUG"
	LogLevelVerbose = "VERBOSE"
	LogLevelInfo    = "INFO"
	LogLevelWarn    = "WARN"
	LogLevelError   = "ERROR"

	JobSequenceStrategyNodeFirst = "node-first"
	JobSequenceStrategyStepFirst = "step-first"
)

type ImportJobInput struct {
	Jobs []*JobInput
}

type JobInput struct {
	UUID               string
	Group              string
	Name               string
	Description        string
	LogLevel           string
	ExecutionEnabled   bool
	MultipleExecutions bool
	NodeFilterEditable bool
	Schedule           *ScheduleInput
	Sequence           *JobSequenceInput
}

type ScheduleInput struct {
	Month   string
	Hour    int
	Minute  int
	Second  int
	Weekday string
	Year    string
}

type JobSequenceInput struct {
	Commands  []string
	KeepGoing bool
	Strategy  string
}

func (i *JobInput) validate() error {
	if i.Name == "" {
		return &Error{
			Code:    ErrCodeInvalidRequest,
			Message: "name is required",
		}
	}
	return nil
}

func newJobDefinition(input *JobInput) (string, error) {
	if err := input.validate(); err != nil {
		return "", fmt.Errorf(": %w", err)
	}
	def := &JobDefinition{
		DefaultTab:         "nodes",
		ExecutionEnabled:   input.ExecutionEnabled,
		Group:              input.Group,
		ID:                 input.UUID,
		UUID:               input.UUID,
		LogLevel:           input.LogLevel,
		MultipleExecutions: input.MultipleExecutions,
		Name:               input.Name,
		Description:        input.Description,
		NodeFilterEditable: input.NodeFilterEditable,
	}

	if input.Sequence != nil {
		cmds := make([]*JobCommand, 0, len(input.Sequence.Commands))
		for _, cmd := range input.Sequence.Commands {
			cmds = append(cmds, &JobCommand{Exec: cmd})
		}
		def.Sequence = &JobSequence{
			Commands:  cmds,
			KeepGoing: input.Sequence.KeepGoing,
			Strategy:  input.Sequence.Strategy,
		}
	}

	if input.Schedule != nil {
		def.Schedule = &JobSchedule{
			Month: input.Schedule.Month,
			Time: &JobScheduleTime{
				Hour:    strconv.Itoa(input.Schedule.Hour),
				Minute:  fmt.Sprintf("%02d", input.Schedule.Minute),
				Seconds: strconv.Itoa(input.Schedule.Second),
			},
			Weekday: &JobScheduleWeekday{Day: input.Schedule.Weekday},
			Year:    input.Schedule.Year,
		}
		def.ScheduleEnabled = true
	}

	b, err := yaml.Marshal([]*JobDefinition{def})
	if err != nil {
		return "", fmt.Errorf(": %w", err)
	}

	return string(b), nil
}

type ImportJobOutput struct {
	Succeeded []*JobDetail `json:"succeeded"`
	Failed    []*JobDetail `json:"failed"`
	Skipped   []*JobDetail `json:"skipped"`
}

type JobDetail struct {
	Index     int    `json:"index"`
	Href      string `json:"href"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Group     string `json:"group"`
	Project   string `json:"project"`
	Permalink string `json:"permalink"`
}

func (c *Client) ImportJob(ctx context.Context, input *ImportJobInput) (*ImportJobOutput, error) {
	var def string
	for _, j := range input.Jobs {
		def2, err := newJobDefinition(j)
		if err != nil {
			return nil, fmt.Errorf(": %w", &Error{
				Code:    ErrCodeUnexpected,
				Message: err.Error(),
			})
		}
		def += fmt.Sprintf("\n%s", def2)
	}

	// Path
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/api/%d/project/%s/jobs/import", c.apiVersion, c.project))

	req2, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader([]byte(def)))
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}

	res, err := c.doRequest(ctx, req2, func(req *http.Request) {
		req.Header.Set("Content-type", "application/yaml")
	})
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var output ImportJobOutput
		if err = json.NewDecoder(res.Body).DecodeContext(ctx, &output); err != nil {
			return nil, fmt.Errorf(": %w", &Error{
				Code:    ErrCodeUnexpected,
				Message: err.Error(),
			})
		}
		return &output, nil
	}

	var v errorRes
	if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	return nil, fmt.Errorf(": %w", v.toError())
}
