package main

import (
	"fmt"
	"regexp"
)

const SplitTagItemsRegexp = `[^;|,"]+|"[^"]*"`

type NestedTagParser struct {
	options []string
	params  map[string]string
}

func Parse(content string) NestedTagParser {
	re := regexp.MustCompile(SplitTagItemsRegexp)
	items := re.FindAllString(content, -1) 

	fmt.Printf("%#v\n", items)

	return NestedTagParser{
		options: []string{},
		params:  map[string]string{},
	}
}

func main() {
	Parse("not null;default:'42'")
}

