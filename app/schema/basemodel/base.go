package basemodel

type BaseResponse struct {
	Data any `json:"data"`
}

type BaseListResponse struct {
	Total int64 `json:"total"`
	Data  any   `json:"data"`
}

type BaseQuery struct {
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}

func NewDefaultQuery() BaseQuery {
	return BaseQuery{
		Skip:  0,
		Limit: 10,
	}
}
