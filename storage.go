package fogg

import (
	"errors"
	"fmt"
	"strings"
)

type Storage struct {
	tags map[string]Tag
}

func Parse(tagContent string) (Storage, error) {
	backticks := []string{`"`}
	delimiters := []string{" "}

	const nameValueDelimiter = ":"

	storage := Storage{
		tags: make(map[string]Tag),
	}

	splitTags, err := splitTagItems(tagContent, true, backticks, delimiters, false)
	if err != nil {
		return storage, err
	}

	for _, t := range splitTags {
		pair := strings.SplitN(t, nameValueDelimiter, 2)
		name, value := pair[0], pair[1]

		value, err = unquoteTagContent(name, value)
		if err != nil {
			return storage, err
		}

		tag, err := ParseSubtag(value, true)
		if err != nil {
			return storage, err
		}
		tag.name = name

		if _, exists := storage.tags[name]; !exists {
			storage.tags[name] = tag
		} else {
			return storage, errors.New(fmt.Sprintf(duplicatedTagsErr, name))
		}
	}

	return storage, nil
}

func (storage *Storage) GetTag(name string) *Tag {
	if tag, exists := storage.tags[name]; exists {
		return &tag
	} else {
		return nil
	}
}

func (storage *Storage) HasTag(name string) bool {
	if _, exists := storage.tags[name]; exists {
		return true
	} else {
		return false
	}
}
