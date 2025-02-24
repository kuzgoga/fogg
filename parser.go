package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func parseTagItems(name string, items []string, argsDelimiter string, trimSpaces bool) (Tag, error) {
	const Separator = ":"

	tag := Tag{
		name:    name,
		params:  make(map[string]TagParam),
		options: make([]string, 0),
	}

	for _, item := range items {
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
				return tag, errors.New(fmt.Sprintf(emptyNameTagErr, value))
			}

			if _, keyExist := tag.params[key]; keyExist {
				return tag, errors.New(fmt.Sprintf(duplicatedParamErr, key))
			}

			param := TagParam{
				Name:  key,
				Value: value,
				Args:  strings.Split(value, argsDelimiter),
			}

			tag.params[key] = param
		} else {
			if trimSpaces {
				tag.options = append(tag.options, strings.TrimSpace(item))
			} else {
				tag.options = append(tag.options, item)
			}
		}
	}
	return tag, nil
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

func splitTagItems(content string, trimSpaces bool, backticks []string, delimiters []string, deleteEscapedSymbols bool) ([]string, error) {
	const EscapeBackslash string = `\`

	var (
		items              []string
		backticksStack     []byte
		currentItem        string
		isBackticksContent bool
	)

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
					items = append(items, strings.TrimSpace(currentItem))
				} else {
					items = append(items, currentItem)
				}
			}
			currentItem = ""
		} else {
			currentItem += string(char)
		}
	}

	if currentItem != "" {
		if trimSpaces {
			items = append(items, strings.TrimSpace(currentItem))
		} else {
			items = append(items, currentItem)
		}
	}

	for i, item := range items {
		for _, quote := range backticks {
			items[i] = strings.ReplaceAll(item, "\\"+quote, quote)
		}
		for _, delimiter := range delimiters {
			items[i] = strings.ReplaceAll(item, "\\"+delimiter, delimiter)
		}
	}

	if len(backticksStack)%2 != 0 {
		return items, errors.New(unclosedBacktickErr)
	}

	return items, nil
}

func unquoteTagContent(name, content string) (string, error) {
	const backtick rune = '"'
	if len(content) < 2 {
		return "", errors.New(fmt.Sprintf(valueLessTagErr, name))
	}
	var startChar, endChar = content[0], content[len(content)-1]
	if rune(startChar) == backtick && rune(endChar) == backtick {
		return content[1 : len(content)-1], nil
	} else {
		return "", errors.New(fmt.Sprintf(nonQuotedValueErr, name))
	}
}

func ParseSubtag(value string, trimSpaces bool) (Tag, error) {
	const argsDelimiter = ","
	subtagBackticks := []string{`'`, `"`}
	subtagDelimiters := []string{";"}

	tagItems, err := splitTagItems(value, trimSpaces, subtagBackticks, subtagDelimiters, true)
	if err != nil {
		return Tag{}, err
	}

	tag, err := parseTagItems("", tagItems, argsDelimiter, true)
	if err != nil {
		return Tag{}, err
	}
	return tag, nil
}
