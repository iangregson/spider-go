package parser_test

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iangregson/spider-go/parser"
)

var emptyHtml = strings.NewReader("")
var invalidHtml = strings.NewReader(`
I am not HTML.
I have no tags.

** Maybe I'm Markdown? **
`)

const indexHtmlPath = "../fixtures/html/index.html"
const blogHtmlPath = "../fixtures/html/blog/index.html"
const baseUrl = "http://localhost:3000"

func TestEmpty(t *testing.T) {
	u, _ := url.Parse(baseUrl)
	wanted := 0
	got := parser.ParseLinks(emptyHtml, u)

	if len(got) != wanted {
		t.Errorf("Expected %d links, got %d", wanted, len(got))
	}
}

func TestInvalid(t *testing.T) {
	u, _ := url.Parse(baseUrl)
	wanted := 0
	got := parser.ParseLinks(invalidHtml, u)

	if len(got) != wanted {
		t.Errorf("Expected %d links, got %d", wanted, len(got))
	}
}

func TestIndex(t *testing.T) {
	filePath, err := filepath.Abs(indexHtmlPath)
	if err != nil {
		t.Errorf("Error getting absolute path: %v\n", err)
		return
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error reading file: %v\n", err)
		return
	}

	u, _ := url.Parse(baseUrl)
	wanted := 11
	got := parser.ParseLinks(bytes.NewReader(content), u)
	if len(got) != wanted {
		t.Errorf("Expected %d links, got %d", wanted, len(got))
	}
}

func TestBlog(t *testing.T) {
	filePath, err := filepath.Abs(blogHtmlPath)
	if err != nil {
		t.Errorf("Error getting absolute path: %v\n", err)
		return
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error reading file: %v\n", err)
		return
	}

	u, _ := url.Parse(baseUrl)
	wanted := 5
	got := parser.ParseLinks(bytes.NewReader(content), u)
	if len(got) != wanted {
		t.Errorf("Expected %d links, got %d", wanted, len(got))
	}
}
