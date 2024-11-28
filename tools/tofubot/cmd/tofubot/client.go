package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	ghToken    string
	httpClient *http.Client
	repo       string
}

func (c *client) createComment(issueNumber int, body string) error {
	requestBody := struct {
		Body string `json:"body"`
	}{body}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%d/comments",
		c.repo,
		issueNumber,
	)
	return httpPost(requestBody, url, c)
}

func httpPost[TReq any](requestBody TReq, url string, c *client) error {
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to encode request body (%w)", err)
	}

	print(url)
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request object (%w)", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.ghToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request to %s (%w)", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to send request to %s (invalid status code: %d)", url, resp.StatusCode)
	}
	return nil
}

func (c *client) addLabels(issueNumber int, labels []string) error {
	return nil
}

func (c *client) getLabels(issueNumber int) ([]string, error) {
	return nil, nil
}

func (c *client) removeLabel(issueNumber int, name string) error {
	return nil
}

func (c *client) addAssignees(issueNumber int, assignees []string) error {
	return nil
}

func (c *client) removeAssignees(issueNumber int, assignees []string) error {
	return nil
}

func (c *client) createLabel(name string, color string, description string) (Label, error) {
	return Label{}, nil
}

func (c *client) listLabels() ([]Label, error) {
	return nil, nil
}

type listPullRequestParams struct {
	State *StateFilter
	Head  string
	Base  string
}

func (c *client) listPullRequests(params listPullRequestParams) ([]PullRequest, error) {
	return nil, nil
}

func (c *client) createPullRequest(title string, body string, head string, base string) (PullRequest, error) {
	return PullRequest{}, nil
}
