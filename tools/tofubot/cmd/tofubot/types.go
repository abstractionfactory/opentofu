package main

type Action string

const ActionCreated Action = "created"
const ActionEdited Action = "edited"
const ActionDeleted Action = "deleted"

type User struct {
	Login string `json:"login"`
}

type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type IssueState string

const IssueStateOpen IssueState = "open"
const IssueStateClosed IssueState = "closed"

type AuthorAssociation string

const AuthorAssociationCollaborator AuthorAssociation = "COLLABORATOR"
const AuthorAssociationContributor AuthorAssociation = "CONTRIBUTOR"
const AuthorAssociationFirstTimer AuthorAssociation = "FIRST_TIMER"
const AuthorAssociationFirstTimeContributor AuthorAssociation = "FIRST_TIME_CONTRIBUTOR"
const AuthorAssociationMannequin AuthorAssociation = "MANNEQUIN"
const AuthorAssociationMember AuthorAssociation = "MEMBER"
const AuthorAssociationNone AuthorAssociation = "NONE"
const AuthorAssociationOwner AuthorAssociation = "OWNER"

type Issue struct {
	Number            int               `json:"number"`
	Title             string            `json:"title"`
	Body              string            `json:"body"`
	User              User              `json:"user"`
	Labels            []Label           `json:"labels"`
	Assignees         []User            `json:"assignee"`
	State             IssueState        `json:"state"`
	AuthorAssociation AuthorAssociation `json:"author_association"`
	PullRequest       *PullRequest      `json:"pull_request,omitempty"`
}

type Team struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type Reference struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	User  *User  `json:"user"`
}

type State string

const StateOpen State = "open"
const StateClosed State = "closed"

type StateFilter string

const StateOpenFilter State = StateOpen
const StateClosedFilter State = StateClosed
const StateFilterAll StateFilter = "all"

type PullRequest struct {
	Issue

	RequestedReviewers []User    `json:"requested_reviewers"`
	RequestedTeams     []Team    `json:"requested_teams"`
	Head               Reference `json:"head"`
	Base               Reference `json:"base"`
	Merged             bool      `json:"merged"`
	MergedBy           *User     `json:"merged_by"`
}

type IssueComment struct {
	ID                int               `json:"id"`
	Body              string            `json:"body"`
	User              User              `json:"user"`
	AuthorAssociation AuthorAssociation `json:"author_association"`
}

type IssueCommentEvent struct {
	Action  Action `json:"action"`
	Changes struct {
		Body *string `json:"body,omitempty"`
	} `json:"changes"`
	Issue   Issue        `json:"issue"`
	Comment IssueComment `json:"comment"`
}

type IssuesEventAction string

const IssuesEventActionOpened IssuesEventAction = "opened"
const IssuesEventActionEdited IssuesEventAction = "edited"
const IssuesEventActionClosed IssuesEventAction = "closed"
const IssuesEventActionReopened IssuesEventAction = "reopened"
const IssuesEventActionAssigned IssuesEventAction = "assigned"
const IssuesEventActionUnassigned IssuesEventAction = "unassigned"
const IssuesEventActionLabeled IssuesEventAction = "labeled"
const IssuesEventActionUnlabeled IssuesEventAction = "unlabeled"

type IssuesEvent struct {
	Action  IssuesEventAction `json:"action"`
	Issue   Issue             `json:"issue"`
	Changes struct {
		Title *string `json:"title,omitempty"`
		Body  *string `json:"body,omitempty"`
	} `json:"changes"`
	Assignee *User  `json:"assignee"`
	Label    *Label `json:"label"`
}
