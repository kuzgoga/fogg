package main

import (
	"slices"
	"testing"
)

func TestSplitTagNormalItems(t *testing.T) {
	parser := NestedTagParser{}

	// Normal tag with options and params
	content := `not null;default:"one,two";check:,name <> 'jinzhu'`
	expected := []string{"not null", "default:\"one,two\"", "check:,name <> 'jinzhu'"}

	err := parser.splitTagItems(content, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expected) {
		t.Errorf("expected %v, got %v", expected, parser.items)
	}
}

func TestSplitTagWithEscape(t *testing.T) {
	parser := NestedTagParser{}

	content := `foreignKey:hy\;pe;index:,unique;`
	expected := []string{"foreignKey:hy;pe", "index:,unique"}

	err := parser.splitTagItems(content, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expected) {
		t.Errorf("expected %#v, got %#v", expected, parser.items)
	}
}

func TestSplitTagWithDoubledDelimiter(t *testing.T) {
	parser := NestedTagParser{}

	content := `not null;;;references:Id;`
	expected := []string{"not null", "references:Id"}

	err := parser.splitTagItems(content, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expected) {
		t.Errorf("expected %#v, got %#v", expected, parser.items)
	}
}

func TestSplitTagWithUnclosedBacktick(t *testing.T) {
	parser := NestedTagParser{}

	content := `not null;references:'Id;`

	err := parser.splitTagItems(content, true)
	if err == nil {
		t.Errorf("expected error: unclosed backtick in tag")
	}
}

func TestParseFunction(t *testing.T) {
	tagContent := `not null;default:'one';check:', n > 1'`
	expectedItems := []string{"not null", "default:'one'", "check:', n > 1'"}

	parser, err := Parse(tagContent)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expectedItems) {
		t.Errorf("expected %v, got %v", expectedItems, parser.items)
	}
}

func TestParseFunctionWithError(t *testing.T) {
	const TagContent = `not null;default:'one;check:', n > 1'`
	const ExpectedError = "unclosed backtick in tag"

	_, err := Parse(TagContent)
	if err != nil && err.Error() != ExpectedError {
		t.Errorf("unexpected error: %s", err)
	} else if err == nil {
		t.Errorf("expected error: %s", ExpectedError)
	}
}
