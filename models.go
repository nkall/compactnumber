package compactnumber

// CompactType is an enum used to specify compaction settings.
type CompactType string

const (
	// No compaction, e.g. 17,000,000.
	None = CompactType("None")
	// Long-form compaction, e.g. 17 million.
	Long = CompactType("Long")
	// Short-form compaction, e.g. 17M.
	Short = CompactType("Short")
)

// CompactForms is a map of compaction type to a slice of form rules, ordered by rule type ascending.
type CompactForms map[CompactType][]CompactFormRule

// CompactFormRule is a single rule for compacting numbers greater than or equal to the value of "type".
// ZeroesInPattern indicates the number of zero characters in the default "other" pattern, and PatternsByPluralForm
// is a map of plural forms (e.g. one, few, many, other, etc.) to patterns.
type CompactFormRule struct {
	Type                 int64
	ZeroesInPattern      int
	PatternsByPluralForm map[string]string
}
