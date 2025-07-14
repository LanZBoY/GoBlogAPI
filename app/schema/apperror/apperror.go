package apperror

import "wentee/blog/app/schema/apperror/errcode"

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (appError *AppError) GetMessage() string {
	return errcode.Message(appError.Code)
}

func New(status int, code string, messge string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: messge,
		Status:  status,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) RawErr() error {
	return e.Err
}
