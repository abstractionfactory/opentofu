package main

import (
	"fmt"
	"regexp"
	"strings"
)

type handler struct {
	client *client
}

var versionRe = regexp.MustCompile(`^v[0-9].[0-9]$`)

func (h *handler) onIssueComment(event IssueCommentEvent) error {
	parts := strings.Split(event.Comment.Body, " ")
	if parts[0] != "/tofubot" {
		return nil
	}

	if len(parts) == 1 || parts[1] == "help" {
		return h.help(event.Issue.Number, event.Comment.User.Login)
	}

	switch parts[1] {
	case "backport":
		if len(parts) != 3 {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤐 Hey @%s, looks like your command doesn't have the right amount of parameters.\n\nTry this: `/tofubot backport vX.Y`",
					event.Comment.User.Login,
				))
		}

		switch event.Comment.AuthorAssociation {
		case AuthorAssociationCollaborator:
		case AuthorAssociationOwner:
		default:
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤐 Hey @%s, you don't seem to have permissions to create a backport. Please ask a maintainer to do this for you. It's nothing personal, I swear.",
					event.Comment.User.Login,
				))
		}

		if event.Issue.PullRequest == nil {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"😵‍💫 Hey @%s, looks like you are asking me to backport an issue? How am I supposed to do that? Try backporting a PR instead.",
					event.Comment.User.Login,
				))
		}
		versionBranch := parts[2]
		if !versionRe.MatchString(versionBranch) {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"😵‍💫 Hey @%s, looks like you have the version branch number mixed up.\n\nTry this: `/tofubot backport vX.Y`",
					event.Comment.User.Login,
				))
		}
		labels, err := h.client.getRepoLabels()
		if err != nil {
			_ = h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤯 Hey @%s, I'm sorry but something went wrong while checking the project labels. Here's the error message:\n\n```\n%s\n```",
					event.Comment.User.Login,
					err,
				))
			return fmt.Errorf("failed to fetch repo labels (%w)", err)
		}
		labelName := "backport/" + versionBranch
		found := false
		for _, label := range labels {
			if label.Name == labelName {
				found = true
				break
			}
		}
		if !found {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"😵‍💫 Hey @%s, looks like the backport label `%s` doesn't exist. Can you check your version number?",
					event.Comment.User.Login,
					labelName,
				))
		}
		if _, err := h.client.addLabels(event.Issue.Number, []string{labelName}); err != nil {
			_ = h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤯 Hey @%s, I'm sorry but something went wrong while adding the backport label. Here's the error message:\n\n```\n%s\n```",
					event.Comment.User.Login,
					err,
				))
			return fmt.Errorf("failed to add labels (%w)", err)
		}
		if event.Issue.State == IssueStateOpen {
			if err = h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"✅ Hey @%s, all right, I've queued up a backport to version %s when the PR merges.",
					event.Comment.User.Login,
					versionBranch,
				)); err != nil {
				return fmt.Errorf("failed to post success comment (%w)", err)
			}
		} else {
			if err = h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"✅ Hey @%s, all right, I'll create a PR to version %s shortly.",
					event.Comment.User.Login,
					versionBranch,
				)); err != nil {
				return fmt.Errorf("failed to post success comment (%w)", err)
			}
		}
		return nil
	case "cancel-backport":
		if len(parts) != 3 {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤐 Hey @%s, looks like your command doesn't have the right amount of parameters.\n\nUsage: `/tofubot cancel-backport vX.Y`",
					event.Comment.User.Login,
				))
		}

		switch event.Comment.AuthorAssociation {
		case AuthorAssociationCollaborator:
		case AuthorAssociationOwner:
		default:
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"🤐 Hey @%s, you don't seem to have permissions to cancel a backport. Please ask a maintainer to do this for you. It's nothing personal, I swear.",
					event.Comment.User.Login,
				))
		}
		if event.Issue.PullRequest == nil {
			return h.client.createComment(
				event.Issue.Number,
				fmt.Sprintf(
					"😵‍💫 Hey @%s, looks like you are asking me to cancel a backport for an issue? How am I supposed to do that? Try backporting a PR instead.",
					event.Comment.User.Login,
				))
		}
	}

	return nil
}

func (h *handler) onIssues(event IssuesEvent) error {
	return nil
}

var help = `👋 Hi @%s! I'm 🤖 TofuBot. I help with all kinds of tasks related to OpenTofu.

<details><summary>

## Here are all the commands I know

</summary>

### ` + "`" + `/tofubot help` + "`" + `

Prints this help text.

### ` + "`" + `/tofubot backport vX.Y` + "`" + `

Backports the current pull request to version vX.Y. If the pull request is not yet merged, it will create the backport PR when the original PR is merged. The current PR must be against the main branch.

### ` + "`" + `/tofubot cancel-backport vX.Y` + "`" + `

Cancels a pending backport to version vX.Y.

</details>
`

func (h *handler) help(issueNumber int, requester string) error {
	return h.client.createComment(issueNumber, fmt.Sprintf(help, requester))
}

func (h *handler) onPullRequest(event PullRequestEvent) error {
	switch event.Action {
	case PullRequestEventActionClosed:
		if !event.PullRequest.Merged {
			return nil
		}
		var pendingLabels []string
		for _, label := range event.PullRequest.Labels {
			if strings.HasPrefix(label.Name, "backport/") && !strings.HasSuffix(label.Name, "/open") && !strings.HasSuffix(label.Name, "/failed") {
				pendingLabels = append(pendingLabels, strings.TrimPrefix(label.Name, "backport/"))
			}
		}
		for _, ver := range pendingLabels {
			if err := createBackport(event.Number, ver, h.client); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func createBackport(number int, ver string, c *client) error {
	return nil
}

func (h *handler) onCreate(event CreateEvent) error {
	fmt.Sprintf("%v", event)
	switch event.RefType {
	case RefTypeBranch:
		ver := strings.TrimPrefix(event.Ref, "branch/")
		if !versionRe.MatchString(ver) {
			return nil
		}
		if _, err := h.client.createLabel("backport/"+ver, "10315b", "Pending backport to version "+ver); err != nil {
			return err
		}
		if _, err := h.client.createLabel("backport/"+ver+"/open", "105b14", "Automatic backport PR to version "+ver+" open"); err != nil {
			return err
		}
		if _, err := h.client.createLabel("backport/"+ver+"/failed", "5b1010", "Automatic backport PR to version "+ver+" failed"); err != nil {
			return err
		}
		return nil
	}
	return nil
}
