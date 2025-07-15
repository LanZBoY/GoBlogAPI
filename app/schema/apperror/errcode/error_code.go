package errcode

const (
	TPYE_ASSERTION_ERROR = "TYPE_ASSETION_ERROR"
	USER_EXIST           = "USER_EXIST"
	USER_NOT_FOUND       = "USER_NOT_FOUND"
	INVALID_TOKEN        = "INVALID_TOKEN"
)

var code2Message = map[string]string{
	USER_EXIST:           "使用者已存在",
	USER_NOT_FOUND:       "使用者不存在",
	INVALID_TOKEN:        "Token 驗證失敗",
	TPYE_ASSERTION_ERROR: "內部推斷型態錯誤",
}

func Message(code string) string {

	if message, ok := code2Message[code]; ok {
		return message
	}

	return "未定義錯誤"
}
