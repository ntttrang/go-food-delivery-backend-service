package datatype

type AppResponse struct {
	Data   any `json:"data"`
	Paging any `json:"paging,omitempty"`
	Filter any `json:"filter,omitempty"`
}

func ResponseSuccess(data any) *AppResponse {
	return &AppResponse{
		Data: data,
	}
}
