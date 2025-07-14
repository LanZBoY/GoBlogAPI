package errcode

const (
	USER_EXIST     = "USER_EXIST"
	USER_NOT_FOUND = "USER_NOT_FOUND"
)

var code2Message = map[string]string{
	USER_EXIST:     "使用者已存在",
	USER_NOT_FOUND: "使用者不存在",
}

func Message(code string) string {

	if message, ok := code2Message[code]; ok {
		return message
	}

	return "未知錯誤"
}
