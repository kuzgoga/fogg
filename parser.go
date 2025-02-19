package main

import (
	"errors"
	"slices"
	"strings"
)

type NestedTagParser struct {
	options []string
	params  map[string]string
	items   []string
}

func (parser *NestedTagParser) splitTagItems(content string, trimSpaces bool) error {
	const EscapeBackSlash = `\`
	backticks := []string{`'`}
	delimiters := []string{`;`}

	currentItem := ""
	var backticksStack []byte
	isBackticksContent := false

	// TODO: process escaped backslashes

	for pos, char := range content {
		priorBackslash := pos != 0 && string(content[pos-1]) == EscapeBackSlash
		if slices.Contains(backticks, string(char)) && !priorBackslash {
			isBacktickStackNotEmpty := len(backticksStack) != 0
			if isBacktickStackNotEmpty && backticksStack[len(backticksStack)-1] == byte(char) {
				isBackticksContent = false
			} else {
				isBackticksContent = true
			}
			backticksStack = append(backticksStack, byte(char))
		}
		if slices.Contains(delimiters, string(char)) && !isBackticksContent && !priorBackslash {
			if currentItem != "" {
				if trimSpaces {
					parser.items = append(parser.items, strings.TrimSpace(currentItem))
				} else {
					parser.items = append(parser.items, currentItem)
				}
			}
			currentItem = ""
		} else {
			currentItem += string(char)
		}
	}

	if currentItem != "" {
		if trimSpaces {
			parser.items = append(parser.items, strings.TrimSpace(currentItem))
		} else {
			parser.items = append(parser.items, currentItem)
		}
	}

	for i, item := range parser.items {
		for _, quote := range backticks {
			parser.items[i] = strings.ReplaceAll(item, "\\"+quote, quote)
		}
		for _, delimiter := range delimiters {
			parser.items[i] = strings.ReplaceAll(item, "\\"+delimiter, delimiter)
		}
	}

	if len(backticksStack)%2 != 0 {
		return errors.New("unclosed backtick in tag")
	}

	return nil
}

func Parse(tagContent string) (NestedTagParser, error) {
	parser := NestedTagParser{}
	if err := parser.splitTagItems(tagContent, true); err != nil {
		return parser, err
	}
	return parser, nil
}
