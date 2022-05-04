package rundeck

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/goccy/go-json"
)

type GetExecutionOutput struct {
	JobExecutionID int    `json:"id"`
	Href           string `json:"href"`
	Permalink      string `json:"permalink"`
	Status         string `json:"status"`
	Project        string `json:"project"`
	ExecutionType  string `json:"executionType"`
	User           string `json:"user"`
	DateStarted    struct {
		UnixTime int64     `json:"unixtime"`
		Date     time.Time `json:"date"`
	} `json:"date-started"`
	JobDetail struct {
		ID              string `json:"id"`
		AverageDuration int    `json:"averageDuration"`
		Name            string `json:"name"`
		Group           string `json:"group"`
		Project         string `json:"project"`
		Description     string `json:"description"`
		Href            string `json:"href"`
		Permalink       string `json:"permalink"`
	} `json:"job"`
	Arg         *string `json:"argString"`
	Description string  `json:"description"`
}

func (c *Client) GetExecution(ctx context.Context, jobExecID int) (*GetExecutionOutput, error) {
	// Path
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/api/%d/execution/%d", c.apiVersion, jobExecID))

	req2, err := http.NewRequest(http.MethodGet, u.String(), nil)
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
		var v GetExecutionOutput
		if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
			return nil, fmt.Errorf(": %w", err)
		}
		return &v, nil
	}

	var v errorRes
	if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	return nil, fmt.Errorf(": %w", v.toError())
}

func (c *Client) DeleteExecution(ctx context.Context, jobExecID int) error {
	// Path
	u, err := url.Parse(c.basePath)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/api/%d/execution/%d", c.apiVersion, jobExecID))

	req2, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}

	res, err := c.doRequest(ctx, req2)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	var v errorRes
	if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	return fmt.Errorf(": %w", v.toError())
}

func (c *Client) AbortExecution(ctx context.Context, jobExecID int) error {
	// Path
	u, err := url.Parse(c.basePath)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/api/%d/execution/%d/abort", c.apiVersion, jobExecID))

	req2, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}

	res, err := c.doRequest(ctx, req2)
	if err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	var v errorRes
	if err = json.NewDecoder(res.Body).DecodeContext(ctx, &v); err != nil {
		return fmt.Errorf(": %w", &Error{
			Code:    ErrCodeUnexpected,
			Message: err.Error(),
		})
	}
	return fmt.Errorf(": %w", v.toError())
}
