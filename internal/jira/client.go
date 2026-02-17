package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/novshi-tech/atl-cli/internal/adf"
	"github.com/novshi-tech/atl-cli/internal/auth"
)

// Client is an HTTP client for the Jira REST API v3.
type Client struct {
	baseURL    string
	email      string
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Jira client from credentials.
func NewClient(creds auth.SiteCredentials) *Client {
	return &Client{
		baseURL:    strings.TrimRight(creds.BaseURL, "/"),
		email:      creds.Email,
		apiToken:   creds.APIToken,
		httpClient: &http.Client{},
	}
}

// NewClientFromStore creates a new Jira client by loading credentials for the given site alias.
func NewClientFromStore(store auth.CredentialStore, siteAlias string) (*Client, error) {
	creds, err := auth.LoadSite(store, siteAlias)
	if err != nil {
		return nil, err
	}
	return NewClient(creds), nil
}

// BaseURL returns the base URL of the Jira instance.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// CreateIssue creates a new issue.
func (c *Client) CreateIssue(project, issueType, summary, description string) (*CreateIssueResponse, error) {
	req := CreateIssueRequest{
		Fields: CreateIssueFields{
			Project:   ProjectKey{Key: project},
			Summary:   summary,
			IssueType: IssueType{Name: issueType},
		},
	}
	if description != "" {
		desc := adf.TextToADF(description)
		req.Fields.Description = &desc
	}

	var resp CreateIssueResponse
	if err := c.doRequest("POST", "/rest/api/3/issue", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateIssue updates an existing issue's summary and/or description.
func (c *Client) UpdateIssue(key, summary, description string) error {
	fields := UpdateIssueFields{}
	if summary != "" {
		fields.Summary = summary
	}
	if description != "" {
		desc := adf.TextToADF(description)
		fields.Description = &desc
	}
	req := UpdateIssueRequest{Fields: fields}
	return c.doRequest("PUT", "/rest/api/3/issue/"+key, req, nil)
}

// TransitionIssue transitions an issue to the given target status name.
func (c *Client) TransitionIssue(key, targetStatus string) error {
	var transResp TransitionsResponse
	if err := c.doRequest("GET", "/rest/api/3/issue/"+key+"/transitions", nil, &transResp); err != nil {
		return fmt.Errorf("fetching transitions: %w", err)
	}

	target := strings.ToLower(targetStatus)
	for _, t := range transResp.Transitions {
		if strings.ToLower(t.Name) == target || strings.ToLower(t.To.Name) == target {
			req := TransitionRequest{Transition: TransitionID{ID: t.ID}}
			return c.doRequest("POST", "/rest/api/3/issue/"+key+"/transitions", req, nil)
		}
	}

	available := make([]string, 0, len(transResp.Transitions))
	for _, t := range transResp.Transitions {
		available = append(available, fmt.Sprintf("%s (→ %s)", t.Name, t.To.Name))
	}
	return fmt.Errorf("no transition matching %q found; available: %s", targetStatus, strings.Join(available, ", "))
}

// AddComment adds a comment to an issue.
func (c *Client) AddComment(key, body string) error {
	req := AddCommentRequest{Body: adf.TextToADF(body)}
	return c.doRequest("POST", "/rest/api/3/issue/"+key+"/comment", req, nil)
}

// SearchIssues searches for issues using JQL.
func (c *Client) SearchIssues(jql string, maxResults int) (*SearchResponse, error) {
	path := fmt.Sprintf("/rest/api/3/search/jql?jql=%s&maxResults=%d&fields=summary,status,issuetype,assignee",
		urlEncode(jql), maxResults)
	var resp SearchResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetIssue retrieves a single issue with full details.
func (c *Client) GetIssue(key string) (*Issue, error) {
	path := fmt.Sprintf("/rest/api/3/issue/%s?fields=summary,status,issuetype,assignee,description,comment", key)
	var resp Issue
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListSprints lists sprints for a board.
func (c *Client) ListSprints(boardID int, state string) (*SprintsResponse, error) {
	path := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint", boardID)
	if state != "" {
		path += "?state=" + urlEncode(state)
	}
	var resp SprintsResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSprintIssues retrieves issues in a sprint.
func (c *Client) GetSprintIssues(sprintID int) (*SprintIssuesResponse, error) {
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue?fields=summary,status,issuetype,assignee", sprintID)
	var resp SprintIssuesResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.String() != "" {
			return fmt.Errorf("jira API error (%d): %s", resp.StatusCode, apiErr.String())
		}
		return fmt.Errorf("jira API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}
	return nil
}
