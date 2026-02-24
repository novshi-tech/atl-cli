package cmd

type JSONIssueItem struct {
	Key      string `json:"key"`
	Summary  string `json:"summary"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Assignee string `json:"assignee"`
}

type JSONIssueDetail struct {
	Key         string            `json:"key"`
	Summary     string            `json:"summary"`
	Status      string            `json:"status"`
	Type        string            `json:"type"`
	Assignee    string            `json:"assignee"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	DueDate     string            `json:"duedate,omitempty"`
	Comments    []JSONCommentItem `json:"comments,omitempty"`
}

type JSONCommentItem struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Body    string `json:"body"`
}

type JSONSprintItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
	Goal  string `json:"goal,omitempty"`
}

type JSONSiteItem struct {
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type JSONMutationResult struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type JSONUserItem struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Active       bool   `json:"active"`
}

type JSONRepoItem struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Language  string `json:"language"`
	IsPrivate bool   `json:"is_private"`
}

type JSONProjectItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type JSONPRItem struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Author string `json:"author"`
	Source string `json:"source"`
	Dest   string `json:"destination"`
}

type JSONRepoDetail struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	IsPrivate   bool   `json:"is_private"`
	MainBranch  string `json:"mainbranch"`
	UpdatedOn   string `json:"updated_on"`
}
