package exponea

// BulkCommand encapsulates a single API command,
// e.g. Event or Customer
type BulkCommand struct {
	Name string `json:"name"`
	Data interface{} `json:"data"`
}

// Bulk is a set of Commands
type Bulk struct {
	Commands []*BulkCommand `json:"commands"`
}

// BulkResponseResult is a single result from the
// bulk request
type BulkResponseResult struct {
	Status    string `json:"status"`
	OtherData string `json:"other_data"`
}

// BulkResponse is a response from the bulk API
type BulkResponse struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
	Success   bool `json:"success"`
	Results   []BulkResponseResult`json:"results"`
}
