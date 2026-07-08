package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
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

	// storyPointsFieldID caches the resolved custom field id for Story
	// Points within the lifetime of this client, since resolving it
	// requires a lookup against /rest/api/3/field.
	storyPointsFieldID    string
	storyPointsFieldKnown bool
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

// CreateIssue creates a new issue. Pass parentKey to set the parent issue
// (an epic for standard issues, or the parent task for sub-tasks). issueType
// may be either an issue type name (e.g. "Task", "タスク") or a numeric id;
// names are resolved against the project's createmeta so that the correct
// project-scoped id is sent — sending only the name can fail with
// "Invalid issue type" when the same name exists in multiple schemes.
func (c *Client) CreateIssue(project, issueType, summary, description, dueDate, parentKey string) (*CreateIssueResponse, error) {
	it, err := c.resolveIssueType(project, issueType)
	if err != nil {
		return nil, err
	}
	req := CreateIssueRequest{
		Fields: CreateIssueFields{
			Project:   ProjectKey{Key: project},
			Summary:   summary,
			IssueType: it,
		},
	}
	if description != "" {
		desc := adf.TextToADF(description)
		req.Fields.Description = &desc
	}
	if dueDate != "" {
		req.Fields.DueDate = dueDate
	}
	if parentKey != "" {
		req.Fields.Parent = &ParentRef{Key: parentKey}
	}

	var resp CreateIssueResponse
	if err := c.doRequest("POST", "/rest/api/3/issue", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// resolveIssueType turns a user-supplied issue type (numeric id or name) into
// an IssueType that the create endpoint can resolve unambiguously. Names are
// looked up via the project's createmeta so we send the project-scoped id.
func (c *Client) resolveIssueType(project, issueType string) (IssueType, error) {
	if issueType == "" {
		return IssueType{}, fmt.Errorf("issue type is required")
	}
	if isAllDigits(issueType) {
		return IssueType{ID: issueType}, nil
	}
	types, err := c.GetIssueTypes(project)
	if err != nil {
		return IssueType{}, fmt.Errorf("resolving issue type %q: %w", issueType, err)
	}
	target := strings.ToLower(issueType)
	for _, t := range types {
		if strings.ToLower(t.Name) == target {
			return IssueType{ID: t.ID}, nil
		}
	}
	names := make([]string, 0, len(types))
	for _, t := range types {
		names = append(names, t.Name)
	}
	return IssueType{}, fmt.Errorf("issue type %q not creatable in project %q; available: %s",
		issueType, project, strings.Join(names, ", "))
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// UpdateIssue updates an existing issue's summary, description, due date, and/or parent.
// parentKey may be an epic key (for standard issues) or a parent task key (for sub-tasks).
func (c *Client) UpdateIssue(key, summary, description, dueDate, parentKey string) error {
	fields := UpdateIssueFields{}
	if summary != "" {
		fields.Summary = summary
	}
	if description != "" {
		desc := adf.TextToADF(description)
		fields.Description = &desc
	}
	if dueDate != "" {
		fields.DueDate = dueDate
	}
	if parentKey != "" {
		fields.Parent = &ParentRef{Key: parentKey}
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

// AssignIssue assigns an issue to a user by accountId.
// Pass nil to unassign.
func (c *Client) AssignIssue(key string, accountID *string) error {
	req := AssignIssueRequest{AccountID: accountID}
	return c.doRequest("PUT", "/rest/api/3/issue/"+key+"/assignee", req, nil)
}

// AddComment adds a comment to an issue.
func (c *Client) AddComment(key, body string) error {
	req := AddCommentRequest{Body: adf.TextToADF(body)}
	return c.doRequest("POST", "/rest/api/3/issue/"+key+"/comment", req, nil)
}

// SearchIssues searches for issues using JQL.
func (c *Client) SearchIssues(jql string, maxResults int) (*SearchResponse, error) {
	spFieldID, err := c.resolveStoryPointsFieldID()
	if err != nil {
		return nil, err
	}
	fieldsParam := "summary,status,issuetype,assignee,parent"
	if spFieldID != "" {
		fieldsParam += "," + spFieldID
	}
	path := fmt.Sprintf("/rest/api/3/search/jql?jql=%s&maxResults=%d&fields=%s",
		urlEncode(jql), maxResults, fieldsParam)
	raw, err := c.doRequestRaw("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp SearchResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}
	populateStoryPoints(raw, resp.Issues, spFieldID)
	return &resp, nil
}

// GetIssue retrieves a single issue with full details.
func (c *Client) GetIssue(key string) (*Issue, error) {
	spFieldID, err := c.resolveStoryPointsFieldID()
	if err != nil {
		return nil, err
	}
	fieldsParam := "summary,status,issuetype,assignee,description,comment,duedate,attachment,parent"
	if spFieldID != "" {
		fieldsParam += "," + spFieldID
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s?fields=%s", key, fieldsParam)
	raw, err := c.doRequestRaw("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Issue
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}
	if spFieldID != "" {
		resp.Fields.StoryPoints = extractFloatField(raw, spFieldID)
	}
	return &resp, nil
}

// GetAttachments returns the attachments for an issue.
func (c *Client) GetAttachments(key string) ([]Attachment, error) {
	path := fmt.Sprintf("/rest/api/3/issue/%s?fields=attachment", key)
	var resp Issue
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Fields.Attachment, nil
}

// DownloadAttachment downloads the content of an attachment by ID and writes it to w.
// Returns the filename reported by the server (from Content-Disposition) if available.
func (c *Client) DownloadAttachment(id string, w io.Writer) (string, error) {
	url := c.baseURL + "/rest/api/3/attachment/content/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("jira API error (%d): %s", resp.StatusCode, string(body))
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return "", fmt.Errorf("writing attachment: %w", err)
	}

	filename := ""
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			filename = params["filename"]
		}
	}
	return filename, nil
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
	spFieldID, err := c.resolveStoryPointsFieldID()
	if err != nil {
		return nil, err
	}
	fieldsParam := "summary,status,issuetype,assignee"
	if spFieldID != "" {
		fieldsParam += "," + spFieldID
	}
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue?fields=%s", sprintID, fieldsParam)
	raw, err := c.doRequestRaw("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp SprintIssuesResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}
	populateStoryPoints(raw, resp.Issues, spFieldID)
	return &resp, nil
}

// GetFields returns all fields known to this Jira site, including custom
// fields. Used to resolve site-specific custom field ids (e.g. Story
// Points) by name.
func (c *Client) GetFields() ([]Field, error) {
	var resp []Field
	if err := c.doRequest("GET", "/rest/api/3/field", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// storyPointsFieldNames lists the field names Jira Cloud uses for story
// points across different project templates: "Story Points" on
// company-managed (classic) software projects, "Story point estimate" on
// team-managed projects.
var storyPointsFieldNames = []string{"story points", "story point estimate"}

// resolveStoryPointsFieldID looks up the custom field id for Story Points on
// this site, caching the result for the lifetime of the client. Returns ""
// (no error) if the site has no such field, so read paths can simply skip
// enrichment; SetStoryPoints treats an empty id as a hard error.
func (c *Client) resolveStoryPointsFieldID() (string, error) {
	if c.storyPointsFieldKnown {
		return c.storyPointsFieldID, nil
	}
	fields, err := c.GetFields()
	if err != nil {
		return "", fmt.Errorf("resolving story points field: %w", err)
	}
	for _, f := range fields {
		for _, candidate := range storyPointsFieldNames {
			if strings.ToLower(f.Name) == candidate {
				c.storyPointsFieldID = f.ID
				c.storyPointsFieldKnown = true
				return f.ID, nil
			}
		}
	}
	c.storyPointsFieldKnown = true
	return "", nil
}

// SetStoryPoints sets the Story Points value on an issue.
func (c *Client) SetStoryPoints(key string, points float64) error {
	fieldID, err := c.resolveStoryPointsFieldID()
	if err != nil {
		return err
	}
	if fieldID == "" {
		return fmt.Errorf("no Story Points field found on this Jira site (looked for: %s)", strings.Join(storyPointsFieldNames, ", "))
	}
	req := map[string]interface{}{
		"fields": map[string]interface{}{fieldID: points},
	}
	return c.doRequest("PUT", "/rest/api/3/issue/"+key, req, nil)
}

// extractFloatField reads a numeric value at raw.fields[fieldID] from a raw
// single-issue JSON response, returning nil if absent, null, or malformed.
// Used to pull dynamically-keyed custom fields (like Story Points) out of
// responses whose static struct fields can't have a per-site JSON tag.
func extractFloatField(raw []byte, fieldID string) *float64 {
	var wrapper struct {
		Fields map[string]json.RawMessage `json:"fields"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil
	}
	v, ok := wrapper.Fields[fieldID]
	if !ok || string(v) == "null" {
		return nil
	}
	var f float64
	if err := json.Unmarshal(v, &f); err != nil {
		return nil
	}
	return &f
}

// populateStoryPoints fills in StoryPoints on each issue from the raw
// "issues" array in a search/sprint-issues response, matching by index.
func populateStoryPoints(raw []byte, issues []Issue, fieldID string) {
	if fieldID == "" {
		return
	}
	var wrapper struct {
		Issues []struct {
			Fields map[string]json.RawMessage `json:"fields"`
		} `json:"issues"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return
	}
	for i := range issues {
		if i >= len(wrapper.Issues) {
			break
		}
		v, ok := wrapper.Issues[i].Fields[fieldID]
		if !ok || string(v) == "null" {
			continue
		}
		var f float64
		if json.Unmarshal(v, &f) == nil {
			issues[i].Fields.StoryPoints = &f
		}
	}
}

// ListProjects lists projects visible to the authenticated user.
func (c *Client) ListProjects(query string, maxResults int) (*ProjectSearchResponse, error) {
	path := fmt.Sprintf("/rest/api/3/project/search?maxResults=%d&orderBy=name", maxResults)
	if query != "" {
		path += "&query=" + urlEncode(query)
	}
	var resp ProjectSearchResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetIssueTypes returns issue types. If projectIDOrKey is non-empty, returns
// only types that the current user can actually create in that project (based
// on the createmeta endpoint, which respects issue type scheme, screen scheme,
// workflow, and permissions). Otherwise returns all globally defined issue
// types.
func (c *Client) GetIssueTypes(projectIDOrKey string) ([]IssueTypeDetail, error) {
	if projectIDOrKey != "" {
		path := "/rest/api/3/issue/createmeta/" + urlEncode(projectIDOrKey) + "/issuetypes"
		var resp CreateMetaIssueTypesResponse
		if err := c.doRequest("GET", path, nil, &resp); err != nil {
			return nil, err
		}
		return resp.IssueTypes, nil
	}
	var resp []IssueTypeDetail
	if err := c.doRequest("GET", "/rest/api/3/issuetype", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMyself returns the currently authenticated user.
func (c *Client) GetMyself() (*User, error) {
	var resp User
	if err := c.doRequest("GET", "/rest/api/3/myself", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchUsers searches for users by display name or email address.
func (c *Client) SearchUsers(query string, maxResults int) ([]User, error) {
	path := fmt.Sprintf("/rest/api/3/user/search?query=%s&maxResults=%d",
		urlEncode(query), maxResults)
	var resp []User
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	respBody, err := c.doRequestRaw(method, path, body)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}
	return nil
}

// doRequestRaw performs the request and returns the raw response body,
// letting callers decode it more than once (e.g. into a typed struct and
// also into a map to pick out a dynamically-keyed custom field).
func (c *Client) doRequestRaw(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.String() != "" {
			return nil, fmt.Errorf("jira API error (%d): %s", resp.StatusCode, apiErr.String())
		}
		return nil, fmt.Errorf("jira API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
