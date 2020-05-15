package errors

import "errors"

var (
	//统一的业务记录不存在的error
	//对于get单条业务记录的接口，在service和storage层都可以返回这个error，只要在interface方法定义里加注释说明
	ErrBusinessRecordNotFound   = errors.New("business record not found")
	ErrBusinessRecordDuplicated = errors.New("business record duplicated")
	ErrBusinessRecordChanged    = errors.New("business record changed")

	ErrIllegalArgument = errors.New("illegal argument")
)
