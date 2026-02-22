package cmd

import (
	"net/http"
	"testing"

	"google.golang.org/api/docs/v1"
	gapi "google.golang.org/api/googleapi"
)

func TestDocsWebViewLink(t *testing.T) {
	if docsWebViewLink("") != "" {
		t.Fatalf("expected empty link")
	}
	link := docsWebViewLink("abc")
	if link != "https://docs.google.com/document/d/abc/edit" {
		t.Fatalf("unexpected link: %q", link)
	}
}

func TestDocsPlainText(t *testing.T) {
	doc := &docs.Document{
		Body: &docs.Body{
			Content: []*docs.StructuralElement{
				{
					Paragraph: &docs.Paragraph{
						Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Hello "}}, {TextRun: &docs.TextRun{Content: "World"}}},
					},
				},
				{
					Table: &docs.Table{
						TableRows: []*docs.TableRow{
							{
								TableCells: []*docs.TableCell{
									{Content: []*docs.StructuralElement{{Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "A"}}}}}}},
									{Content: []*docs.StructuralElement{{Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "B"}}}}}}},
								},
							},
						},
					},
				},
			},
		},
	}

	text := docsPlainText(doc, 0)
	if text == "" {
		t.Fatalf("expected text output")
	}
	if text != "Hello WorldA\tB" {
		t.Fatalf("unexpected docs text: %q", text)
	}

	limited := docsPlainText(doc, 5)
	if limited != "Hello" {
		t.Fatalf("unexpected limited text: %q", limited)
	}
}

func TestBodyEndIndex(t *testing.T) {
	if bodyEndIndex(nil) != 0 {
		t.Fatal("expected 0 for nil body")
	}
	if bodyEndIndex(&docs.Body{}) != 0 {
		t.Fatal("expected 0 for empty body")
	}
	body := &docs.Body{
		Content: []*docs.StructuralElement{
			{EndIndex: 1},
			{EndIndex: 50},
		},
	}
	if got := bodyEndIndex(body); got != 49 {
		t.Fatalf("expected 49, got %d", got)
	}
}

func TestSetRequestsTabId(t *testing.T) {
	requests := []*docs.Request{
		{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: &docs.Range{StartIndex: 1, EndIndex: 10},
			},
		},
		{
			UpdateTextStyle: &docs.UpdateTextStyleRequest{
				Range: &docs.Range{StartIndex: 1, EndIndex: 5},
			},
		},
		{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{Index: 1},
				Text:     "hello",
			},
		},
		{
			InsertText: &docs.InsertTextRequest{
				EndOfSegmentLocation: &docs.EndOfSegmentLocation{},
				Text:                 "world",
			},
		},
		{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{StartIndex: 1, EndIndex: 20},
			},
		},
		{
			InsertTable: &docs.InsertTableRequest{
				Location: &docs.Location{Index: 5},
				Rows:     2,
				Columns:  2,
			},
		},
	}

	setRequestsTabId(requests, "t.abc")

	if requests[0].UpdateParagraphStyle.Range.TabId != "t.abc" {
		t.Fatalf("UpdateParagraphStyle Range TabId not set")
	}
	if requests[1].UpdateTextStyle.Range.TabId != "t.abc" {
		t.Fatalf("UpdateTextStyle Range TabId not set")
	}
	if requests[2].InsertText.Location.TabId != "t.abc" {
		t.Fatalf("InsertText Location TabId not set")
	}
	if requests[3].InsertText.EndOfSegmentLocation.TabId != "t.abc" {
		t.Fatalf("InsertText EndOfSegmentLocation TabId not set")
	}
	if requests[4].DeleteContentRange.Range.TabId != "t.abc" {
		t.Fatalf("DeleteContentRange Range TabId not set")
	}
	if requests[5].InsertTable.Location.TabId != "t.abc" {
		t.Fatalf("InsertTable Location TabId not set")
	}
}

func TestSetRequestsTabId_EmptyIsNoop(t *testing.T) {
	requests := []*docs.Request{
		{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: &docs.Range{StartIndex: 1, EndIndex: 10},
			},
		},
	}

	setRequestsTabId(requests, "")

	if requests[0].UpdateParagraphStyle.Range.TabId != "" {
		t.Fatal("expected TabId to remain empty")
	}
}

func TestIsDocsNotFound(t *testing.T) {
	if isDocsNotFound(&gapi.Error{Code: http.StatusNotFound}) != true {
		t.Fatalf("expected not found")
	}
	if isDocsNotFound(&gapi.Error{Code: http.StatusForbidden}) {
		t.Fatalf("unexpected not found")
	}
}
