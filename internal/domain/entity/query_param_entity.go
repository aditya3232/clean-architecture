package entity

type QueryStringEntity struct {
	Search    string
	Page      int64
	Limit     int64
	OrderBy   string
	OrderType string
}
