package jira

import "novshi-tech.com/jira-cli/internal/adf"

// CreateIssueRequest is the request body for creating an issue.
type CreateIssueRequest struct {
	Fields CreateIssueFields `json:"fields"`
}

type CreateIssueFields struct {
	Project     ProjectKey `json:"project"`
	Summary     string     `json:"summary"`
	IssueType   IssueType  `json:"issuetype"`
	Description *adf.Node  `json:"description,omitempty"`
}

type ProjectKey struct {
	Key string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

// CreateIssueResponse is the response from creating an issue.
type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// UpdateIssueRequest is the request body for updating an issue.
type UpdateIssueRequest struct {
	Fields UpdateIssueFields `json:"fields"`
}

type UpdateIssueFields struct {
	Summary     string    `json:"summary,omitempty"`
	Description *adf.Node `json:"description,omitempty"`
}

// AddCommentRequest is the request body for adding a comment.
type AddCommentRequest struct {
	Body adf.Node `json:"body"`
}

// TransitionsResponse is the response from the transitions endpoint.
type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	ID   string           `json:"id"`
	Name string           `json:"name"`
	To   TransitionStatus `json:"to"`
}

type TransitionStatus struct {
	Name string `json:"name"`
}

// TransitionRequest is the request body for transitioning an issue.
type TransitionRequest struct {
	Transition TransitionID `json:"transition"`
}

type TransitionID struct {
	ID string `json:"id"`
}

// APIError represents an error response from the Jira API.
type APIError struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

func (e *APIError) String() string {
	msg := ""
	for _, m := range e.ErrorMessages {
		msg += m + "; "
	}
	for k, v := range e.Errors {
		msg += k + ": " + v + "; "
	}
	return msg
}
