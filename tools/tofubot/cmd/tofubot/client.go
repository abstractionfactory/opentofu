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
	type reqType struct {
		Body string `json:"body"`
	}
	requestBody := reqType{body}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%d/comments",
		c.repo,
		issueNumber,
	)
	_, err := httpPost[reqType, any](requestBody, url, c)
	return err
}

func (c *client) addLabels(issueNumber int, labels []string) ([]Label, error) {
	type reqType struct {
		Labels []string `json:"labels"`
	}
	requestBody := reqType{labels}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/issues/%s/labels",
		c.repo,
		issueNumber,
	)
	return httpPost[reqType, []Label](requestBody, url, c)
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
	type reqType struct {
		Name        string `json:"name"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}
	requestBody := reqType{name, color, description}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/labels",
		c.repo,
	)
	return httpPost[reqType, Label](requestBody, url, c)
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

func (c *client) getRepoLabels() ([]Label, error) {
	// TODO add pagination
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/labels?per_page=100",
		c.repo,
	)
	return httpGet[[]Label](url, c)
}

func httpGet[TResp any](url string, c *client) (value TResp, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return value, fmt.Errorf("failed to create HTTP request object (%w)", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.ghToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return value, fmt.Errorf("send request to %s (%w)", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode > 299 {
		return value, fmt.Errorf("failed to send request to %s (invalid status code: %d)", url, resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("failed to decode Github response (%w)", err)
	}
	return value, nil
}

func httpPost[TReq any, TResp any](requestBody TReq, url string, c *client) (value TResp, err error) {
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return value, fmt.Errorf("failed to encode request body (%w)", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return value, fmt.Errorf("failed to create HTTP request object (%w)", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.ghToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return value, fmt.Errorf("send request to %s (%w)", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode > 299 {
		return value, fmt.Errorf("failed to send request to %s (invalid status code: %d)", url, resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("failed to decode Github response (%w)", err)
	}
	return value, nil
}
