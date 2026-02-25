package bitbucket

// Repository represents a Bitbucket repository.
type Repository struct {
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Description string     `json:"description"`
	IsPrivate   bool       `json:"is_private"`
	Language    string     `json:"language"`
	UpdatedOn   string     `json:"updated_on"`
	MainBranch  *Branch    `json:"mainbranch"`
	Links       RepoLinks  `json:"links"`
}

type Branch struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type RepoLinks struct {
	HTML Link `json:"html"`
}

type Link struct {
	Href string `json:"href"`
}

// RepositoriesResponse is the paginated response for listing repositories.
type RepositoriesResponse struct {
	PageLen int          `json:"pagelen"`
	Size    int          `json:"size"`
	Page    int          `json:"page"`
	Values  []Repository `json:"values"`
}

// PullRequestsResponse is the paginated response for listing pull requests.
type PullRequestsResponse struct {
	PageLen int           `json:"pagelen"`
	Size    int           `json:"size"`
	Page    int           `json:"page"`
	Values  []PullRequest `json:"values"`
}

// PullRequest represents a Bitbucket pull request.
type PullRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	Description string    `json:"description"`
	Source      PRRef     `json:"source"`
	Destination PRRef     `json:"destination"`
	Author      PRUser    `json:"author"`
	CreatedOn   string    `json:"created_on"`
	UpdatedOn   string    `json:"updated_on"`
	Links       PRLinks   `json:"links"`
}

type PRRef struct {
	Branch   Branch     `json:"branch"`
	Repository *PRRepo `json:"repository"`
}

type PRRepo struct {
	FullName string `json:"full_name"`
}

type PRUser struct {
	DisplayName string `json:"display_name"`
	Nickname    string `json:"nickname"`
}

type PRLinks struct {
	HTML Link `json:"html"`
}

// CreatePRRequest is the request body for creating a pull request.
type CreatePRRequest struct {
	Title       string       `json:"title"`
	Source      CreatePRRef  `json:"source"`
	Destination *CreatePRRef `json:"destination,omitempty"`
	Description string       `json:"description,omitempty"`
}

type CreatePRRef struct {
	Branch CreatePRBranch `json:"branch"`
}

type CreatePRBranch struct {
	Name string `json:"name"`
}

// PRComment represents a comment on a pull request.
type PRComment struct {
	ID        int              `json:"id"`
	Content   PRCommentContent `json:"content"`
	User      PRUser           `json:"user"`
	CreatedOn string           `json:"created_on"`
	UpdatedOn string           `json:"updated_on"`
	Inline    *PRInline        `json:"inline"`
}

type PRCommentContent struct {
	Raw    string `json:"raw"`
	Markup string `json:"markup"`
	HTML   string `json:"html"`
}

type PRInline struct {
	Path string `json:"path"`
	From *int   `json:"from"`
	To   *int   `json:"to"`
}

// PRCommentsResponse is the paginated response for listing PR comments.
type PRCommentsResponse struct {
	PageLen int         `json:"pagelen"`
	Size    int         `json:"size"`
	Page    int         `json:"page"`
	Values  []PRComment `json:"values"`
	Next    string      `json:"next"`
}

// BBUser represents the authenticated Bitbucket user.
type BBUser struct {
	DisplayName string `json:"display_name"`
	Nickname    string `json:"nickname"`
	AccountID   string `json:"account_id"`
	UUID        string `json:"uuid"`
	CreatedOn   string `json:"created_on"`
}

// APIError represents an error response from the Bitbucket API.
type APIError struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Message string `json:"message"`
}

func (e *APIError) String() string {
	return e.Error.Message
}
