package bizerr

import (
	"errors"
	"fmt"
	"github.com/xhigher/hzgo/resp"
)

type ErrType string

type Error struct {
	Type ErrType
	Err error
}

const (
	typeNotFound = "not_found"
	typeAlreadyExists = "already_exists"
	typePasswordWrong = "password_wrong"

	typeUserNull = "user_null"
	typeUserExists = "user_exists"
	typeUserBlocked = "user_blocked"
	typeUserCanceled = "user_canceled"
)

var (
	errTypeMsgMap = map[ErrType]string{
		typeNotFound:"记录不存在",
		typePasswordWrong:"密码错误",
		typeUserNull:"用户不存在",
		typeUserExists:"用户已存在",
		typeUserBlocked:"用户已封禁",
		typeUserCanceled:"用户已注销",
	}
)

func New(err error) *Error{
	return &Error{
		Err:err,
	}
}

func newErr(typ ErrType, msg string) error {
	if len(msg) == 0 {
		msg = errTypeMsgMap[typ]
	}
	if len(msg) == 0 {
		return nil
	}
	return errors.New(msg)
}

func NotFound(msg string) *Error{
	return &Error{
		Type:typeNotFound,
		Err: newErr(typeNotFound, msg),
	}
}

func AlreadyExists(msg string) *Error{
	return &Error{
		Type:typeAlreadyExists,
		Err: newErr(typeAlreadyExists, msg),
	}
}

func PasswordWrong(msg string) *Error{
	return &Error{
		Type:typePasswordWrong,
		Err: newErr(typePasswordWrong, msg),
	}
}

func UserNull() *Error{
	return &Error{
		Type:typePasswordWrong,
		Err: newErr(typePasswordWrong, ""),
	}
}

func UserExists() *Error{
	return &Error{
		Type:typeUserExists,
		Err: newErr(typeUserExists, ""),
	}
}

func UserBlocked() *Error{
	return &Error{
		Type:typeUserBlocked,
		Err: newErr(typeUserBlocked, ""),
	}
}

func UserCanceled() *Error{
	return &Error{
		Type:typeUserCanceled,
		Err: newErr(typeUserCanceled, ""),
	}
}

func (e *Error) String() string{
	return fmt.Sprintf("{type:%s, err:%v}", e.Type, e.Err)
}

func (e *Error) ToResp() resp.BaseResp{
	if len(e.Type) == 0 {
		if e.Err == nil {
			return resp.ErrorInternal
		}
		return resp.BaseResp{
			Code: resp.ErrorInternal.Code,
			Msg: fmt.Sprintf("%s: %s", resp.ErrorInternal.Msg,e.Err.Error()),
		}
	}
	return resp.BaseResp{
		Code: resp.OK.Code,
		Msg:  e.Err.Error(),
	}
}

func (e *Error) NotFound() bool {
	return e.Type == typeNotFound
}

func (e *Error) AlreadyExists() bool {
	return e.Type == typeAlreadyExists
}

func (e *Error) PasswordWrong() bool {
	return e.Type == typePasswordWrong
}

func (e *Error) UserNull() bool {
	return e.Type == typeUserNull
}

func (e *Error) UserExists() bool {
	return e.Type == typeUserExists
}

func (e *Error) UserBlocked() bool {
	return e.Type == typeUserBlocked
}

func (e *Error) UserCanceled() bool {
	return e.Type == typeUserCanceled
}
