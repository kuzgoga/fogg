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

func TestEscapedBackslashes(t *testing.T) {
	tagContent := `references:badId\\;polymorphic:my\'BadValue;`
	expectedItems := []string{`references:badId\`, `polymorphic:my\'BadValue`}

	parser, err := Parse(tagContent)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expectedItems) {
		t.Errorf("expected %v, got %v", expectedItems, parser.items)
	}
}

func TestSplitTagTrimSpaces(t *testing.T) {
	parser := NestedTagParser{}

	content := `  not null  ;  default:'one,two'  ;  check:,name <> 'jinzhu'  `
	expected := []string{"  not null  ", "  default:'one,two'  ", "  check:,name <> 'jinzhu'  "}

	err := parser.splitTagItems(content, false)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(parser.items, expected) {
		t.Errorf("expected %v, got %v", expected, parser.items)
	}
}

func TestFindEscapedBackslashesIndexes(t *testing.T) {
	input := `This is a test string with escaped backslashes \\ and more \\ backslashes`
	expectedEscaped := []int{47, 59}
	expectedIgnored := []int{48, 60}

	escaped, ignored := findEscapedBackslashesIndexes(input)

	if !slices.Equal(escaped, expectedEscaped) {
		t.Errorf("expected escaped indexes %v, got %v", expectedEscaped, escaped)
	}

	if !slices.Equal(ignored, expectedIgnored) {
		t.Errorf("expected ignored indexes %v, got %v", expectedIgnored, ignored)
	}
}

func TestFindEscapedBackslashesIndexesWithBreak(t *testing.T) {
	input := `This is a test string with an escaped backslash at the end \\`
	expectedEscaped := []int{59}
	expectedIgnored := []int{60}

	escaped, ignored := findEscapedBackslashesIndexes(input)

	if !slices.Equal(escaped, expectedEscaped) {
		t.Errorf("expected escaped indexes %v, got %v", expectedEscaped, escaped)
	}

	if !slices.Equal(ignored, expectedIgnored) {
		t.Errorf("expected ignored indexes %v, got %v", expectedIgnored, ignored)
	}
}
