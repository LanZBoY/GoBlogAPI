package errcode

const (
	TPYE_ASSERTION_ERROR = "TYPE_ASSETION_ERROR"
	USER_EXIST           = "USER_EXIST"
	USER_NOT_FOUND       = "USER_NOT_FOUND"
	INVALID_TOKEN        = "INVALID_TOKEN"
	EXPIRED_TOKEN        = "EXPIRED_TOKEN"
	BAD_REQUEST          = "BAD_REQUEST"
)

var code2Message = map[string]string{
	USER_EXIST:           "使用者已存在",
	USER_NOT_FOUND:       "使用者不存在",
	INVALID_TOKEN:        "Token 驗證失敗",
	EXPIRED_TOKEN:        "Token 過期",
	BAD_REQUEST:          "輸入錯誤",
	TPYE_ASSERTION_ERROR: "內部推斷型態錯誤",
}

func Message(code string) string {

	if message, ok := code2Message[code]; ok {
		return message
	}

	return "未定義錯誤"
}
