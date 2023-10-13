package models

type StackPushRequest struct {
	Value int64 `json:"value"`
}

type StackPushRangeRequest struct {
	Values []int64 `json:"values"`
}
