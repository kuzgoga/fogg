package fogg

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	const tag = `gorm:"default:'ui\\path';index:,unique;not null;foreignKey:Customer\"Id"`
	expectedTag := Tag{
		name:  "gorm",
		value: "not null",
		params: map[string]TagParam{
			"default": {
				Name:  "default",
				Value: `ui\path`,
				Args:  []string{`ui\path`},
			},
			"foreignKey": {
				Name:  "foreignKey",
				Value: `Customer"Id`,
				Args:  []string{`Customer"Id`},
			},
			"index": {
				Name:  "index",
				Value: `,unique`,
				Args:  []string{``, `unique`},
			},
		},
		options: []string{
			"not null",
		},
	}

	storage, err := Parse(tag)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !storage.HasTag("gorm") {
		t.Errorf("`gorm` tag is absont in storage")
	}

	parsedTag := storage.GetTag("gorm")

	if !reflect.DeepEqual(expectedTag, *parsedTag) {
		t.Errorf("got: %+v\n, expected: %+v\n", parsedTag, expectedTag)
	}
}

func TestParseTagErrors(t *testing.T) {
	const nonQuotedValueTag string = `gorm:default:value'`
	expectedError := fmt.Sprintf(nonQuotedValueErr, "gorm")
	_, err := Parse(nonQuotedValueTag)
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error `%s`, got: %s", expectedError, err)
	}
}

func TestParseSubtagErrors(t *testing.T) {
	const unquotedTag string = `gorm:"default:'value"`
	const expectedError = unclosedBacktickErr
	_, err := Parse(unquotedTag)
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error `%s`, got: %s", expectedError, err)
	}
}

func TestParseDuplicatedTagsError(t *testing.T) {
	const duplicatedTags string = `gorm:"not null" gorm:"foreignKey:CustomerId"`
	expectedError := fmt.Sprintf(duplicatedTagsErr, "gorm")
	_, err := Parse(duplicatedTags)
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error `%s`, got: %s", expectedError, err)
	}
}

func TestParseSplitTagItemsError(t *testing.T) {
	const invalidTagContent = `gorm:"default:'value`
	expectedError := unclosedBacktickErr
	_, err := Parse(invalidTagContent)
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error `%s`, got: %s", expectedError, err)
	}
}

func TestGetTagNotFound(t *testing.T) {
	storage := Storage{
		tags: make(map[string]Tag),
	}

	tag := storage.GetTag("nonexistent")
	if tag != nil {
		t.Errorf("expected nil, got %v", tag)
	}
}

func TestHasTagNotFound(t *testing.T) {
	storage := Storage{
		tags: make(map[string]Tag),
	}

	if storage.HasTag("nonexistent") {
		t.Errorf("expected false, got true")
	}
}
