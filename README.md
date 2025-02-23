# Nested tags
Parser for complex tags (like gorm tags)

## Escape rule
> Only values in backticks can be escaped

Example:
```go
`gorm:"default:'\"escapedValue\"';" json:"ui\\workspace"`
```
