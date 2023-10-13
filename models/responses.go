package models

type BaseResponse struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

type StackResponse struct {
	BaseResponse
	StackSize uint64 `json:"stack_size"`
}

type StackPopResponse struct {
	StackResponse
	Value int64 `json:"value"`
}

type StackTopResponse struct {
	StackResponse
	Value int64 `json:"value"`
}

type StackPushResponse struct {
	StackResponse
}

type StackPushRangeResponse struct {
	StackResponse
}

type MainResponse struct {
	BaseResponse
	Message string `json:"message"`
}
