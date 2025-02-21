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

	err := parser.splitSubtagItems(content, true)
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

	err := parser.splitSubtagItems(content, true)
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

	err := parser.splitSubtagItems(content, true)
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

	err := parser.splitSubtagItems(content, true)
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
	expectedItems := []string{`references:badId\\`, `polymorphic:my\'BadValue`}

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

	err := parser.splitSubtagItems(content, false)
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

func TestDistributeItemsToOptionsAndParams(t *testing.T) {
	parser := NestedTagParser{
		items:  []string{"option1", "param1:value1", "option2", "param2:value2"},
		params: make(map[string]string),
	}

	err := parser.distributeItemsToOptionsAndParams(true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{"option1", "option2"}
	expectedParams := map[string]string{"param1": "value1", "param2": "value2"}

	if !slices.Equal(parser.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, parser.options)
	}

	for key, value := range expectedParams {
		if parser.params[key] != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, parser.params[key])
		}
	}
}

func TestDistributeItemsToOptionsAndParamsWithSpaces(t *testing.T) {
	parser := NestedTagParser{
		items:  []string{" option1 ", " param1 : value1 ", " option2 ", " param2 : value2 "},
		params: make(map[string]string),
	}

	err := parser.distributeItemsToOptionsAndParams(true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{"option1", "option2"}
	expectedParams := map[string]string{"param1": "value1", "param2": "value2"}

	if !slices.Equal(parser.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, parser.options)
	}

	for key, value := range expectedParams {
		if parser.params[key] != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, parser.params[key])
		}
	}
}

func TestDistributeItemsToOptionsAndParamsWithDuplicateParam(t *testing.T) {
	parser := NestedTagParser{
		items:  []string{"param1:value1", "param1:value2"},
		params: make(map[string]string),
	}

	err := parser.distributeItemsToOptionsAndParams(true)
	if err == nil {
		t.Errorf("expected error for duplicate param, got nil")
	}
}

func TestDistributeItemsToOptionsAndParamsWithEmptyKey(t *testing.T) {
	parser := NestedTagParser{
		items:  []string{"param1:value1", ":value2"},
		params: make(map[string]string),
	}

	err := parser.distributeItemsToOptionsAndParams(true)
	if err == nil {
		t.Errorf("expected error for empty key, got nil")
	}
}

func TestDistributeItemsToOptionsAndParamsWithoutTrimSpaces(t *testing.T) {
	parser := NestedTagParser{
		items:  []string{" option1 ", " param1 : value1 ", " option2 ", " param2 : value2 "},
		params: make(map[string]string),
	}

	err := parser.distributeItemsToOptionsAndParams(false)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{" option1 ", " option2 "}
	expectedParams := map[string]string{" param1 ": " value1 ", " param2 ": " value2 "}

	if !slices.Equal(parser.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, parser.options)
	}

	for key, value := range expectedParams {
		if parser.params[key] != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, parser.params[key])
		}
	}
}

// TODO: process all branches of deleteEscapeSymbols
