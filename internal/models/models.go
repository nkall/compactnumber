package models

// CompactFormRule is a single rule for compacting numbers greater than or equal to the value of "type".
// ZeroesInPattern indicates the number of zero characters in the default "other" pattern, and PatternsByPluralForm
// is a map of plural forms (e.g. one, few, many, other, etc.) to patterns.
type CompactFormRule struct {
	Type                 int64
	ZeroesInPattern      int
	PatternsByPluralForm map[string]string
}
