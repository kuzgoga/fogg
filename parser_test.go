package main

import (
	"slices"
	"testing"
)

const argsDelimiter = `,`

func TestParseFunction(t *testing.T) {
	tagContent := `not null;default:'one';check:', n > 1'`
	expectedItems := []string{"not null", "default:'one'", "check:', n > 1'"}
	backticks := []string{`'`}
	delimiters := []string{`;`}

	items, err := splitTagItems(tagContent, true, backticks, delimiters, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(items, expectedItems) {
		t.Errorf("expected %v, got %v", expectedItems, items)
	}
}

func TestParseFunctionWithError(t *testing.T) {
	const tagContent = `not null;default:'one;check:', n > 1'`
	const expectedError = "unclosed backtick in tag"

	backticks := []string{`'`}
	delimiters := []string{`;`}

	_, err := splitTagItems(tagContent, true, backticks, delimiters, true)
	if err != nil && err.Error() != expectedError {
		t.Errorf("unexpected error: %s", err)
	} else if err == nil {
		t.Errorf("expected error: %s", expectedError)
	}
}

func TestEscapedBackslashes(t *testing.T) {
	tagContent := `references:badId\\;polymorphic:my\'BadValue;`
	expectedItems := []string{`references:badId\\`, `polymorphic:my\'BadValue`}

	backticks := []string{`"`}
	delimiters := []string{`;`}
	items, err := splitTagItems(tagContent, true, backticks, delimiters, false)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(items, expectedItems) {
		t.Errorf("expected %v, got %v", expectedItems, items)
	}
}

func TestSplitTagTrimSpaces(t *testing.T) {
	content := `  not null  ;  default:'one,two'  ;  check:,name <> 'jinzhu'  `
	expected := []string{"  not null  ", "  default:'one,two'  ", "  check:,name <> 'jinzhu'  "}
	backticks := []string{`"`}
	delimiters := []string{`;`}

	items, err := splitTagItems(content, false, backticks, delimiters, false)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}
}

/* --- FindEscapedBackslashes function tests --- */
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

/* --- Distribution function tests --- */
func TestDistributeItemsToOptionsAndParams(t *testing.T) {
	items := []string{"option1", "param1:value1", "option2", "param2:value2"}

	tag, err := parseTagItems("", items, argsDelimiter, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{"option1", "option2"}
	expectedParams := map[string]string{"param1": "value1", "param2": "value2"}

	if !slices.Equal(tag.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, tag.options)
	}

	for key, value := range expectedParams {
		if tag.params[key].Value != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, tag.params[key].Value)
		}
	}
}

func TestDistributeItemsToOptionsAndParamsWithSpaces(t *testing.T) {
	items := []string{" option1 ", " param1 : value1 ", " option2 ", " param2 : value2 "}

	tag, err := parseTagItems("", items, argsDelimiter, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{"option1", "option2"}
	expectedParams := map[string]string{"param1": "value1", "param2": "value2"}

	if !slices.Equal(tag.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, tag.options)
	}

	for key, value := range expectedParams {
		if tag.params[key].Value != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, tag.params[key].Value)
		}
	}
}

func TestDistributeItemsToOptionsAndParamsWithDuplicateParam(t *testing.T) {
	items := []string{"param1:value1", "param1:value2"}

	_, err := parseTagItems("", items, argsDelimiter, true)
	if err == nil {
		t.Errorf("expected error for duplicate param, got nil")
	}
}

func TestDistributeItemsToOptionsAndParamsWithEmptyKey(t *testing.T) {
	items := []string{"param1:value1", ":value2"}

	_, err := parseTagItems("", items, argsDelimiter, true)
	if err == nil {
		t.Errorf("expected error for empty key, got nil")
	}
}

func TestDistributeItemsToOptionsAndParamsWithoutTrimSpaces(t *testing.T) {
	items := []string{" option1 ", " param1 : value1 ", " option2 ", " param2 : value2 "}

	tag, err := parseTagItems("", items, argsDelimiter, false)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedOptions := []string{" option1 ", " option2 "}
	expectedParams := map[string]string{" param1 ": " value1 ", " param2 ": " value2 "}

	if !slices.Equal(tag.options, expectedOptions) {
		t.Errorf("expected options %v, got %v", expectedOptions, tag.options)
	}

	for key, value := range expectedParams {
		if tag.params[key].Value != value {
			t.Errorf("expected param %s to be %s, got %s", key, value, tag.params[key].Value)
		}
	}
}

// TODO: process all branches of deleteEscapeSymbols
