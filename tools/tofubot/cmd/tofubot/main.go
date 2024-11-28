package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func usage() {
	print("Usage: GITHUB_REPOSITORY_OWNER=... GITHUB_REPOSITORY=... GITHUB_EVENT_NAME=... GITHUB_EVENT_PATH=... GITHUB_TOKEN=... go run ./cmd/tofubot\n")
}

func main() {
	repo := os.Getenv("GITHUB_REPOSITORY")
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	githubToken := os.Getenv("GITHUB_TOKEN")

	if repo == "" {
		print("Missing GITHUB_REPOSITORY.")
		usage()
		os.Exit(1)
	}
	if eventName == "" {
		print("Missing GITHUB_EVENT_NAME.")
		usage()
		os.Exit(1)
	}
	if eventPath == "" {
		print("Missing GITHUB_EVENT_PATH.")
		usage()
		os.Exit(1)
	}
	if githubToken == "" {
		print("Missing GITHUB_TOKEN.")
		usage()
		os.Exit(1)
	}

	hnd := &handler{
		client: &client{
			ghToken:    githubToken,
			httpClient: http.DefaultClient,
			repo:       repo,
		},
	}

	eventContents, err := os.ReadFile(eventPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	switch eventName {
	case "issues":
		err = handle[IssuesEvent](eventContents, hnd.onIssues)
	case "issue_comment":
		err = handle[IssueCommentEvent](eventContents, hnd.onIssueComment)
	case "pull_request":
		err = handle[PullRequestEvent](eventContents, hnd.onPullRequest)
	case "create":
		err = handle[CreateEvent](eventContents, hnd.onCreate)
	default:
		log.Fatalf("Unhandled event type: %s", eventName)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func handle[T any](contents []byte, handler func(event T) error) error {
	var ev T

	if err := json.Unmarshal(contents, &ev); err != nil {
		return err
	}

	return handler(ev)
}
