package adf

import "strings"

// Node represents an Atlassian Document Format node.
type Node struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Version int    `json:"version,omitempty"`
	Content []Node `json:"content,omitempty"`
}

// TextToADF converts plain text into an ADF document node.
// Each line becomes a separate paragraph.
func TextToADF(text string) Node {
	lines := strings.Split(text, "\n")
	paragraphs := make([]Node, 0, len(lines))
	for _, line := range lines {
		p := Node{Type: "paragraph"}
		if line != "" {
			p.Content = []Node{{Type: "text", Text: line}}
		}
		paragraphs = append(paragraphs, p)
	}
	return Node{
		Type:    "doc",
		Version: 1,
		Content: paragraphs,
	}
}
