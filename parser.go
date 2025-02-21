package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type NestedTagParser struct {
	name    string
	options []string
	params  map[string]string
	items   []string
}

func (parser *NestedTagParser) distributeItemsToOptionsAndParams(trimSpaces bool) error {
	const Separator = ":"
	for _, item := range parser.items {
		if strings.Contains(item, Separator) {
			var key, value string
			pair := strings.SplitN(item, Separator, 2)

			if trimSpaces {
				key = strings.TrimSpace(pair[0])
				value = strings.TrimSpace(pair[1])
			} else {
				key = pair[0]
				value = pair[1]
			}

			if len(key) == 0 {
				return errors.New(fmt.Sprintf("invalid param with empty name and value \"%s\"", value))
			}

			if _, keyExist := parser.params[key]; keyExist {
				return errors.New(fmt.Sprintf("duplicated param \"%s\" in tag", key))
			}

			parser.params[key] = value
		} else {
			if trimSpaces {
				parser.options = append(parser.options, strings.TrimSpace(item))
			} else {
				parser.options = append(parser.options, item)
			}
		}
	}
	return nil
}

func findEscapedBackslashesIndexes(s string) (escapedBackslashes []int, ignoredBackslashes []int) {
	const EscapedBackslash = `\\`
	var escapedIndexes []int
	var ignoredIndexes []int
	index := strings.Index(s, EscapedBackslash)

	for index != -1 {
		escapedIndexes = append(escapedIndexes, index)
		ignoredIndexes = append(ignoredIndexes, index+1)
		if index+len(EscapedBackslash) < len(s) {
			index = strings.Index(s[index+len(EscapedBackslash):], EscapedBackslash)
			if index != -1 {
				index += escapedIndexes[len(escapedIndexes)-1] + len(EscapedBackslash)
			}
		} else {
			break
		}
	}
	return escapedIndexes, ignoredIndexes
}

func (parser *NestedTagParser) splitSubtagItems(content string, trimSpaces bool) error {
	backticks := []string{`'`}
	delimiters := []string{`;`}
	return parser.splitTagItems(content, trimSpaces, backticks, delimiters, false)
}

func (parser *NestedTagParser) splitPrimaryTagsItems(content string, trimSpaces bool) error {
	backticks := []string{`"`}
	delimiters := []string{` `}
	return parser.splitTagItems(content, trimSpaces, backticks, delimiters, false)
}

func (parser *NestedTagParser) splitTagItems(content string, trimSpaces bool, backticks []string, delimiters []string, deleteEscapedSymbols bool) error {
	const EscapeBackslash = `\`

	currentItem := ""
	var backticksStack []byte
	isBackticksContent := false

	escapedBackslashes, ignoredBackslashes := findEscapedBackslashesIndexes(content)

	for pos, char := range content {
		// Process escaped backslashes
		if slices.Contains(escapedBackslashes, pos) {
			currentItem += string(char)
			continue
		}
		if slices.Contains(ignoredBackslashes, pos) && deleteEscapedSymbols {
			continue
		} else if slices.Contains(ignoredBackslashes, pos) && !deleteEscapedSymbols {
			currentItem += string(char)
			continue
		}

		isPriorBackslashEscaped := pos != 0 && slices.Contains(ignoredBackslashes, pos-1)
		priorBackslash := pos != 0 && string(content[pos-1]) == EscapeBackslash && !isPriorBackslashEscaped

		// Skip backslashes used for escaping symbols
		if priorBackslash && slices.Contains(backticks, string(char)) && deleteEscapedSymbols {
			currentItem = currentItem[:len(currentItem)-1] + string(char)
			continue
		}

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
	parser := NestedTagParser{
		params: map[string]string{},
	}
	if err := parser.splitSubtagItems(tagContent, true); err != nil {
		return parser, err
	}
	if err := parser.distributeItemsToOptionsAndParams(true); err != nil {
		return parser, err
	}
	return parser, nil
}
