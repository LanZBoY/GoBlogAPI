package errcode

import "testing"

func TestMessage(t *testing.T) {
	if msg := Message(USER_NOT_FOUND); msg != "使用者不存在" {
		t.Errorf("unexpected message: %s", msg)
	}

	if msg := Message("UNKNOWN"); msg != "未定義錯誤" {
		t.Errorf("unexpected message for unknown code: %s", msg)
	}
}
