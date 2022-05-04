package rundeck

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type NewClientInput struct {
	BasePath   string
	Token      string
	APIVersion int
	Project    string
}

type Client struct {
	httpClient *http.Client
	basePath   string
	token      string
	apiVersion int
	project    string
}

func NewClient(input *NewClientInput) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		basePath:   input.BasePath,
		token:      input.Token,
		apiVersion: input.APIVersion,
		project:    input.Project,
	}
}

func (c Client) doRequest(ctx context.Context, req *http.Request, opts ...func(*http.Request)) (*http.Response, error) {
	req.Header.Set("X-Rundeck-Auth-Token", c.token)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")

	for _, opt := range opts {
		opt(req)
	}

	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return res, nil
}

type errorRes struct {
	Error      bool   `json:"error"`
	APIVersion int    `json:"apiversion"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

func (r errorRes) toError() *Error {
	// ジョブ実行中の削除エラー
	if r.ErrorCode == "api.error.exec.delete.failed" && strings.Contains(r.Message, "The execution is currently running") {
		return &Error{
			Code:    ErrCodeDeletedRunningExecution,
			Message: r.Message,
		}
	}
	return &Error{
		Code:    ErrCodeUnexpected,
		Message: r.Message,
	}
}
