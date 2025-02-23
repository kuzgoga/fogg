package main

type Tag struct {
	name    string
	params  map[string]TagParam
	options []string
}

type TagParam struct {
	Name  string
	Value string
	Args  []string
}

func ParseFromString(content string) (Tag, error) {
	panic("not implemented")
	return Tag{}, nil
}

func ParseSubtag(content string) (Tag, error) {
	tag := Tag{
		name:    "",
		params:  make(map[string]TagParam),
		options: nil,
	}
	panic("not implemented")
	return tag, nil
}

func (tag *Tag) HasOption(name string) {
	panic("not implemented")
}

func (tag *Tag) GetParam(name string) *TagParam {
	panic("not implemented")
}

func (tag *Tag) Name() string {
	return tag.name
}
