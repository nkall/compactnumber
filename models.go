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
