package fogg

import (
	"slices"
)

// TODO: introduce constants here
// TODO: config feature

type Tag struct {
	name    string
	value   string
	params  map[string]TagParam
	options []string
}

func (tag *Tag) Name() string {
	return tag.name
}

func (tag *Tag) HasOption(name string) bool {
	if slices.Contains(tag.options, name) {
		return true
	} else {
		return false
	}
}

func (tag *Tag) HasParam(name string) bool {
	if _, exist := tag.params[name]; exist {
		return true
	} else {
		return false
	}
}

func (tag *Tag) GetParam(name string) *TagParam {
	if param, exist := tag.params[name]; exist {
		return &param
	} else {
		return nil
	}
}

func (tag *Tag) GetParamOr(name string, defaultValue string) string {
	if param, exist := tag.params[name]; exist {
		return param.Value
	} else {
		return defaultValue
	}
}

func (tag *Tag) GetValue() string {
	return tag.value
}
