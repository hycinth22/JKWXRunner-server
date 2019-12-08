// service包及下属包，保证其error可对客户安全显示。
package service

// 向客户隐藏的内部错误，使用此类型包装
type InternalError struct {
	Detail error
}

func (err InternalError) Error() string {
	return "内部错误"
}

func (err InternalError) Unwrap() error {
	return err.Detail
}

func IsInternalError(err error) bool {
	_, yes := err.(InternalError)
	return yes
}

// hidden the actual error into an InternalError.
// call UnwrapInternalError when the actual error is needed.
func WrapAsInternalError(err error) error {
	return InternalError{Detail: err}
}

// return actual error of an InternalError.
// if the parameter is not an InternalError, return it as it is
func UnwrapInternalError(err error) error {
	interErr, yes := err.(InternalError)
	if !yes {
		return err
	}
	return interErr.Detail
}
