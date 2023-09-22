package parse_test

import (
	"fmt"
	"testing"

	"github.com/sanzhang007/crawl_nodes/parse"
)

type changfengT struct {
	Items []Item
}

type Item struct {
	Name        string
	Path        string
	ContentType string
}

func TestXxx(t *testing.T) {
	fmt.Printf("parse.Getpaths(): %v\n", parse.GetChangfengpaths())
}
