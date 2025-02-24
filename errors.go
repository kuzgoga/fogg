package main

const (
	duplicatedTagsErr   string = `duplicated tags with name "%s"`
	emptyNameTagErr     string = `invalid param with empty name and value "%s"`
	duplicatedParamErr  string = `duplicated param "%s" in tag`
	unclosedBacktickErr string = `unclosed backtick in tag`
	valueLessTagErr     string = "Invalid `%s` tag syntax"
	nonQuotedValueErr   string = "`%s` tag value must be in quotation marks"
)
