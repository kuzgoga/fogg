package fogg

import (
	"fmt"
	"reflect"
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

	backticks := []string{`'`}
	delimiters := []string{`;`}
	items, err := splitTagItems(tagContent, true, backticks, delimiters, false)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !slices.Equal(items, expectedItems) {
		t.Errorf("expected %v, got %v", expectedItems, items)
	}
}

func TestDeleteEscapedBackslashes(t *testing.T) {
	tagContent := `references:badId\\;polymorphic:my\'BadValue;`
	expectedItems := []string{`references:badId\`, `polymorphic:my'BadValue`}

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

func TestSplitTagTrimSpaces(t *testing.T) {
	content := `  not null  ;  default:'one,two'  ;  check:,name <> 'jinzhu'  `
	expected := []string{"  not null  ", "  default:'one,two'  ", "  check:,name <> 'jinzhu'  "}
	backticks := []string{`'`}
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

func TestUnquoteTag(t *testing.T) {
	quoted := `"content"`
	expectedUnquoted := `content`
	invalidQuoted := `"`
	expectedInvalidQuotedError := fmt.Sprintf(valueLessTagErr, "")
	nonQuoted := `value`
	nonQuotedValueError := fmt.Sprintf(nonQuotedValueErr, "")

	if unquoted, err := unquoteTagContent("", quoted); err != nil || unquoted != expectedUnquoted {
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if unquoted != expectedUnquoted {
			t.Errorf("wrong unquoted: got `%s`, expected `%s`", unquoted, expectedUnquoted)
		}
	}

	if _, err := unquoteTagContent("", invalidQuoted); err == nil || err.Error() != expectedInvalidQuotedError {
		t.Errorf("unexpected error: got `%s`, expected `%s`", err, expectedInvalidQuotedError)
	}

	if _, err := unquoteTagContent("", nonQuoted); err == nil || err.Error() != nonQuotedValueError {
		t.Errorf("unexpected error: got `%s`, expected `%s`", err, nonQuotedValueError)
	}
}

func TestParseSubtag(t *testing.T) {
	const validTag string = `default:'\"SomeValue';foreignKey:CustomerId;`

	expectedTag := Tag{
		name: "",
		params: map[string]TagParam{
			"default": {
				Name:  "default",
				Value: `"SomeValue`,
				Args:  []string{`"SomeValue`},
			},
			"foreignKey": {
				Name:  "foreignKey",
				Value: "CustomerId",
				Args:  []string{"CustomerId"},
			},
		},
		options: []string{},
	}

	tag, err := ParseSubtag(validTag, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(expectedTag, tag) {
		t.Errorf("got: %+v\n, expected: %+v\n", tag, expectedTag)
	}
}

func TestParseSubtagsErrors(t *testing.T) {
	const unclosedBacktickTag = `default:'"SomeValue';foreignKey:CustomerId;`
	if _, err := ParseSubtag(unclosedBacktickTag, true); err == nil || err.Error() != unclosedBacktickErr {
		t.Errorf("expected error `%s`, got `%s`", unclosedBacktickErr, err)
	}

	const duplicatedParamsTag = `foreignKey:CustomerId;foreignKey:UserId;`
	duplicatedParamExpectedErr := fmt.Sprintf(duplicatedParamErr, "foreignKey")
	if _, err := ParseSubtag(duplicatedParamsTag, true); err == nil || err.Error() != duplicatedParamExpectedErr {
		t.Errorf("expected error `%s`, got `%s`", duplicatedParamExpectedErr, err)
	}
}

func TestUnquoteValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'value'", "value"},
		{`"value"`, "value"},
		{"value", "value"},
		{"'unmatched", "'unmatched"},
		{`"unmatched`, `"unmatched`},
	}

	for _, test := range tests {
		result := unquoteParamValue(test.input)
		if result != test.expected {
			t.Errorf("unquoteParamValue(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestValueDistribution(t *testing.T) {
	const tag string = `foo;omitempty`
	const expectedValue = `foo`

	parsedTag, err := ParseSubtag(tag, true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if parsedTag.value != expectedValue {
		t.Errorf("unexpected value: `%s` got, expected `%s`", expectedValue, tag)
	}
}
