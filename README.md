[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/kuzgoga/fogg/graph/badge.svg?token=X8N0PNCKXB)](https://codecov.io/gh/kuzgoga/fogg)
# Fogg
Parser for tags of golang structures with built-in validation and support for complex tags (according the GORM tags style).

## Overview
* user-friendly API
* escaping support
* detailed error messages
* support for `GORM` and `classic` tags styles
* high test coverage

## Installation
```
go get -u github.com/kuzgoga/fogg@v0.1.2
```

## Example
```go
package main

import (
	"fmt"
	"github.com/kuzgoga/fogg"
)

func main() {
	const tag = `gorm:"default:'something';not null"`
	tags, err := fogg.Parse(tag)

	if err != nil {
		fmt.Printf("Tag validation error: %s\n", err)
	}

	defaultValue := tags.GetTag("gorm").GetParam("default").Value
	isNotNull := tags.GetTag("gorm").HasOption("not null")

	fmt.Printf("Default value: %s\n", defaultValue) // > Default value: something
	fmt.Printf("Not null: %v\n", isNotNull)         // > Not null: true
}
```

## License
Released under the [MIT License](https://github.com/kuzgoga/fogg/blob/master/LICENSE)
