package jira

import "github.com/novshi-tech/atl-cli/internal/adf"

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

// SearchResponse is the response from the JQL search endpoint.
type SearchResponse struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

// Issue represents a Jira issue.
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string         `json:"summary"`
	Status      Status         `json:"status"`
	IssueType   IssueTypeInfo  `json:"issuetype"`
	Assignee    *User          `json:"assignee"`
	Description *adf.Node      `json:"description"`
	Comment     *CommentResult `json:"comment"`
}

type Status struct {
	Name string `json:"name"`
}

type IssueTypeInfo struct {
	Name string `json:"name"`
}

type User struct {
	DisplayName string `json:"displayName"`
}

type CommentResult struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Author  User     `json:"author"`
	Body    adf.Node `json:"body"`
	Created string   `json:"created"`
}

// SprintsResponse is the response from the sprint list endpoint.
type SprintsResponse struct {
	Values []Sprint `json:"values"`
}

type Sprint struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
	Goal  string `json:"goal,omitempty"`
}

// SprintIssuesResponse is the response from the sprint issues endpoint.
type SprintIssuesResponse struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
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
