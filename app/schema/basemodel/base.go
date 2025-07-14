package basemodel

type BaseResponse struct {
	Data any `json:"data"`
}

type BaseListResponse struct {
	Total int64 `json:"total"`
	Data  any   `json:"data"`
}

type BaseQuery struct {
	Skip  int64 `json:"skip" binding:"default=0"`
	Limit int64 `json:"limit" binding:"default=10"`
}
