package main

import (
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

	parsedTag := storage.GetTag("gorm")

	if !reflect.DeepEqual(expectedTag, *parsedTag) {
		t.Errorf("got: %+v\n, expected: %+v\n", parsedTag, expectedTag)
	}
}
