package victor

// TimeStat represents timing statistics for an operation
type TimeStat struct {
	Count uint64  `json:"count"`
	Total float64 `json:"total"`
	Last  float64 `json:"last"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// IndexStats represents aggregate statistics for the index
type IndexStats struct {
	Insert  TimeStat `json:"insert"`
	Delete  TimeStat `json:"delete"`
	Dump    TimeStat `json:"dump"`
	Search  TimeStat `json:"search"`
	SearchN TimeStat `json:"search_n"`
}
