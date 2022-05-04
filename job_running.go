package rundeck

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/goccy/go-json"
)

type RunJobInput struct {
	JobID   string                 `json:"-"`
	RunAt   *time.Time             `json:"runAtTime,omitempty"`
	Options map[string]interface{} `json:"options"`
}

func (c *Client) RunJob(ctx context.Context, input *RunJobInput) (*GetExecutionOutput, error) {
	if input.RunAt != nil && time.Now().After(*input.RunAt) {
		input.RunAt = nil
	}

	// Path
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/api/%d/job/%s/executions", c.apiVersion, input.JobID))

	b, err := json.MarshalContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}

	req2, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}

	res, err := c.doRequest(ctx, req2)
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var output GetExecutionOutput
		if err = json.NewDecoder(res.Body).DecodeContext(ctx, &output); err != nil {
			return nil, fmt.Errorf(": %w", err)
		}
		return &output, nil
	}

	var v errorRes
	if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	return nil, fmt.Errorf(": %w", v.toError())
}
