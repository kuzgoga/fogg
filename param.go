package main

import "slices"

type TagParam struct {
	Name  string
	Value string
	Args  []string
}

func (param *TagParam) HasArg(arg string) bool {
	if slices.Contains(param.Args, arg) {
		return true
	} else {
		return false
	}
}
