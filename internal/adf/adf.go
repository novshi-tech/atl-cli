package adf

import (
	"regexp"
	"strings"
)

// Node represents an Atlassian Document Format node.
type Node struct {
	Type    string         `json:"type"`
	Text    string         `json:"text,omitempty"`
	Version int            `json:"version,omitempty"`
	Content []Node         `json:"content,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
}

// Mark represents an ADF inline mark (bold, italic, link, etc.).
type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// TextToADF converts markdown-like text into an ADF document node.
// Supported syntax:
//   - ## Heading       → heading (level 2-6)
//   - * item / - item  → bulletList
//   - ||h1||h2||       → table header row
//   - |c1|c2|          → table data row
//   - **bold**         → strong mark
//   - [text](url)      → link mark
//   - blank lines      → empty paragraph
func TextToADF(text string) Node {
	lines := strings.Split(text, "\n")
	var blocks []Node
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Heading: ## text
		if level, content, ok := parseHeading(line); ok {
			blocks = append(blocks, makeHeading(level, content))
			i++
			continue
		}

		// Table: lines starting with || or |
		if isTableLine(line) {
			var tableLines []string
			for i < len(lines) && isTableLine(lines[i]) {
				tableLines = append(tableLines, lines[i])
				i++
			}
			blocks = append(blocks, makeTable(tableLines))
			continue
		}

		// Bullet list: lines starting with * or -
		if isBulletLine(line) {
			var items []string
			for i < len(lines) && isBulletLine(lines[i]) {
				items = append(items, trimBullet(lines[i]))
				i++
			}
			blocks = append(blocks, makeBulletList(items))
			continue
		}

		// Empty line → empty paragraph
		if strings.TrimSpace(line) == "" {
			blocks = append(blocks, Node{Type: "paragraph"})
			i++
			continue
		}

		// Regular paragraph
		blocks = append(blocks, makeParagraph(line))
		i++
	}

	return Node{
		Type:    "doc",
		Version: 1,
		Content: blocks,
	}
}

// --- Heading ---

func parseHeading(line string) (int, string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "#") {
		return 0, "", false
	}
	level := 0
	for _, ch := range trimmed {
		if ch == '#' {
			level++
		} else {
			break
		}
	}
	if level < 1 || level > 6 {
		return 0, "", false
	}
	content := strings.TrimSpace(trimmed[level:])
	if content == "" {
		return 0, "", false
	}
	return level, content, true
}

func makeHeading(level int, text string) Node {
	return Node{
		Type:    "heading",
		Attrs:   map[string]any{"level": level},
		Content: parseInline(text),
	}
}

// --- Bullet list ---

var bulletPrefix = regexp.MustCompile(`^\s*[\*\-]\s+`)

func isBulletLine(line string) bool {
	return bulletPrefix.MatchString(line)
}

func trimBullet(line string) string {
	loc := bulletPrefix.FindStringIndex(line)
	if loc == nil {
		return line
	}
	return line[loc[1]:]
}

func makeBulletList(items []string) Node {
	listItems := make([]Node, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, Node{
			Type: "listItem",
			Content: []Node{
				{Type: "paragraph", Content: parseInline(item)},
			},
		})
	}
	return Node{Type: "bulletList", Content: listItems}
}

// --- Table ---

func isTableLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "|")
}

func makeTable(tableLines []string) Node {
	rows := make([]Node, 0, len(tableLines))
	for _, line := range tableLines {
		trimmed := strings.TrimSpace(line)
		isHeader := strings.HasPrefix(trimmed, "||")
		row := parseTableRow(trimmed, isHeader)
		rows = append(rows, row)
	}
	return Node{Type: "table", Content: rows}
}

func parseTableRow(line string, isHeader bool) Node {
	var cells []string
	if isHeader {
		// ||col1||col2|| format
		line = strings.TrimPrefix(line, "||")
		line = strings.TrimSuffix(line, "||")
		cells = strings.Split(line, "||")
	} else {
		// |col1|col2| format
		line = strings.TrimPrefix(line, "|")
		line = strings.TrimSuffix(line, "|")
		cells = strings.Split(line, "|")
	}

	cellType := "tableCell"
	if isHeader {
		cellType = "tableHeader"
	}

	cellNodes := make([]Node, 0, len(cells))
	for _, cell := range cells {
		cellNodes = append(cellNodes, Node{
			Type: cellType,
			Content: []Node{
				{Type: "paragraph", Content: parseInline(strings.TrimSpace(cell))},
			},
		})
	}

	return Node{Type: "tableRow", Content: cellNodes}
}

// --- Paragraph ---

func makeParagraph(text string) Node {
	return Node{
		Type:    "paragraph",
		Content: parseInline(text),
	}
}

// --- Inline parsing (bold, links) ---

type inlineMatch struct {
	start, end int
	node       Node
}

var (
	boldRe    = regexp.MustCompile(`\*\*(.+?)\*\*`)
	linkRe    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	mentionRe = regexp.MustCompile(`@\[([^\]:]+):([^\]]+)\]`)
)

// parseInline splits text into ADF text nodes with marks for **bold**, [text](url), and @[name:accountId] mentions.
func parseInline(text string) []Node {
	if text == "" {
		return nil
	}

	var matches []inlineMatch

	for _, loc := range boldRe.FindAllStringSubmatchIndex(text, -1) {
		matches = append(matches, inlineMatch{
			start: loc[0],
			end:   loc[1],
			node: Node{
				Type:  "text",
				Text:  text[loc[2]:loc[3]],
				Marks: []Mark{{Type: "strong"}},
			},
		})
	}

	for _, loc := range linkRe.FindAllStringSubmatchIndex(text, -1) {
		matches = append(matches, inlineMatch{
			start: loc[0],
			end:   loc[1],
			node: Node{
				Type:  "text",
				Text:  text[loc[2]:loc[3]],
				Marks: []Mark{{Type: "link", Attrs: map[string]any{"href": text[loc[4]:loc[5]]}}},
			},
		})
	}

	for _, loc := range mentionRe.FindAllStringSubmatchIndex(text, -1) {
		displayName := text[loc[2]:loc[3]]
		accountID := text[loc[4]:loc[5]]
		matches = append(matches, inlineMatch{
			start: loc[0],
			end:   loc[1],
			node: Node{
				Type:  "mention",
				Attrs: map[string]any{"id": accountID, "text": "@" + displayName},
			},
		})
	}

	if len(matches) == 0 {
		return []Node{{Type: "text", Text: text}}
	}

	// Sort matches by start position
	sortMatches(matches)

	// Remove overlapping matches (keep first)
	var filtered []inlineMatch
	lastEnd := 0
	for _, m := range matches {
		if m.start >= lastEnd {
			filtered = append(filtered, m)
			lastEnd = m.end
		}
	}

	// Build result nodes
	var nodes []Node
	pos := 0
	for _, m := range filtered {
		if m.start > pos {
			nodes = append(nodes, Node{Type: "text", Text: text[pos:m.start]})
		}
		nodes = append(nodes, m.node)
		pos = m.end
	}
	if pos < len(text) {
		nodes = append(nodes, Node{Type: "text", Text: text[pos:]})
	}

	return nodes
}

func sortMatches(matches []inlineMatch) {
	for i := 1; i < len(matches); i++ {
		for j := i; j > 0 && matches[j].start < matches[j-1].start; j-- {
			matches[j], matches[j-1] = matches[j-1], matches[j]
		}
	}
}
