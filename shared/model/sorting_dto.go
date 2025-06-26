package sharedmodel

type SortingDto struct {
	SortBy    string `json:"sortBy" form:"sortBy"`
	Direction string `json:"direction" form:"direction"`
}
