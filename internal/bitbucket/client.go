package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/novshi-tech/atl-cli/internal/auth"
)

const baseURL = "https://api.bitbucket.org/2.0"

// Client is an HTTP client for the Bitbucket Cloud REST API 2.0.
type Client struct {
	email      string
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Bitbucket client from credentials.
// If BBAPIToken is set, it is used instead of APIToken for authentication.
func NewClient(creds auth.SiteCredentials) *Client {
	token := creds.APIToken
	if creds.BBAPIToken != "" {
		token = creds.BBAPIToken
	}
	return &Client{
		email:      creds.Email,
		apiToken:   token,
		httpClient: &http.Client{},
	}
}

// NewClientFromStore creates a new Bitbucket client by loading credentials for the given site alias.
func NewClientFromStore(store auth.CredentialStore, siteAlias string) (*Client, error) {
	creds, err := auth.LoadSite(store, siteAlias)
	if err != nil {
		return nil, err
	}
	return NewClient(creds), nil
}

// GetCurrentUser returns the currently authenticated user.
func (c *Client) GetCurrentUser() (*BBUser, error) {
	var resp BBUser
	if err := c.doRequest("GET", "/user", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListRepositories lists repositories for a workspace.
func (c *Client) ListRepositories(workspace string, page, pagelen int) (*RepositoriesResponse, error) {
	path := fmt.Sprintf("/repositories/%s?page=%d&pagelen=%d", workspace, page, pagelen)
	var resp RepositoriesResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRepository retrieves a single repository.
func (c *Client) GetRepository(workspace, repoSlug string) (*Repository, error) {
	path := fmt.Sprintf("/repositories/%s/%s", workspace, repoSlug)
	var resp Repository
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPullRequests lists pull requests for a repository.
func (c *Client) ListPullRequests(workspace, repoSlug, state string, page, pagelen int) (*PullRequestsResponse, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests?page=%d&pagelen=%d", workspace, repoSlug, page, pagelen)
	if state != "" {
		path += "&state=" + state
	}
	var resp PullRequestsResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreatePullRequest creates a new pull request.
func (c *Client) CreatePullRequest(workspace, repoSlug string, req CreatePRRequest) (*PullRequest, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests", workspace, repoSlug)
	var resp PullRequest
	if err := c.doRequest("POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPRComments lists comments on a pull request.
func (c *Client) ListPRComments(workspace, repoSlug string, prID int) (*PRCommentsResponse, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments", workspace, repoSlug, prID)
	var resp PRCommentsResponse
	if err := c.doRequest("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) doRequest(method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
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
			return fmt.Errorf("bitbucket API error (%d): %s", resp.StatusCode, apiErr.String())
		}
		return fmt.Errorf("bitbucket API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}
	return nil
}
