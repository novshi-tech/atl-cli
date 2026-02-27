package adf

import (
	"encoding/json"
	"testing"
)

func TestParseInline_Italic(t *testing.T) {
	nodes := parseInline("hello *world* end")
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "hello ", nil)
	assertTextNode(t, nodes[1], "world", []string{"em"})
	assertTextNode(t, nodes[2], " end", nil)
}

func TestParseInline_BoldNotAffected(t *testing.T) {
	nodes := parseInline("hello **bold** end")
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "hello ", nil)
	assertTextNode(t, nodes[1], "bold", []string{"strong"})
	assertTextNode(t, nodes[2], " end", nil)
}

func TestParseInline_BoldAndItalic(t *testing.T) {
	nodes := parseInline("**bold** and *italic*")
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "bold", []string{"strong"})
	assertTextNode(t, nodes[1], " and ", nil)
	assertTextNode(t, nodes[2], "italic", []string{"em"})
}

func TestParseInline_ItalicAndBold(t *testing.T) {
	nodes := parseInline("*italic* and **bold**")
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "italic", []string{"em"})
	assertTextNode(t, nodes[1], " and ", nil)
	assertTextNode(t, nodes[2], "bold", []string{"strong"})
}

func TestParseInline_OnlyItalic(t *testing.T) {
	nodes := parseInline("*emphasis*")
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "emphasis", []string{"em"})
}

func TestParseInline_MultipleItalics(t *testing.T) {
	nodes := parseInline("*one* plain *two*")
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %s", len(nodes), jsonStr(nodes))
	}
	assertTextNode(t, nodes[0], "one", []string{"em"})
	assertTextNode(t, nodes[1], " plain ", nil)
	assertTextNode(t, nodes[2], "two", []string{"em"})
}

func TestTextToADF_ItalicInParagraph(t *testing.T) {
	doc := TextToADF("This is *important* text")
	if doc.Type != "doc" {
		t.Fatalf("expected doc node, got %s", doc.Type)
	}
	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Content))
	}
	para := doc.Content[0]
	if para.Type != "paragraph" {
		t.Fatalf("expected paragraph, got %s", para.Type)
	}
	if len(para.Content) != 3 {
		t.Fatalf("expected 3 inline nodes, got %d: %s", len(para.Content), jsonStr(para.Content))
	}
	assertTextNode(t, para.Content[1], "important", []string{"em"})
}

// --- helpers ---

func assertTextNode(t *testing.T, n Node, text string, markTypes []string) {
	t.Helper()
	if n.Type != "text" {
		t.Errorf("expected type text, got %s", n.Type)
		return
	}
	if n.Text != text {
		t.Errorf("expected text %q, got %q", text, n.Text)
	}
	if len(markTypes) == 0 {
		if len(n.Marks) != 0 {
			t.Errorf("expected no marks for %q, got %v", text, n.Marks)
		}
		return
	}
	if len(n.Marks) != len(markTypes) {
		t.Errorf("expected %d marks for %q, got %d: %v", len(markTypes), text, len(n.Marks), n.Marks)
		return
	}
	for i, mt := range markTypes {
		if n.Marks[i].Type != mt {
			t.Errorf("expected mark type %q at index %d, got %q", mt, i, n.Marks[i].Type)
		}
	}
}

func jsonStr(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
