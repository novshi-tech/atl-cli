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
