package main

import (
	"fmt"
	"strings"
)

type handler struct {
	client *client
}

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
					"😵‍💫 Hey @%s, looks like you are asking me to backport an issue? How am I supposed to do that?",
					event.Comment.User.Login,
				))
		}

	case "cancel-backport":
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
					"😵‍💫 Hey @%s, looks like you are asking me to cancel a backport for an issue? How am I supposed to do that?",
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
