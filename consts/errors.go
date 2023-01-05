package consts

import "errors"

var (
	ErrUserNotExisted    = errors.New("user not existed")
	ErrUserPasswordWrong = errors.New("user password wrong")
)
