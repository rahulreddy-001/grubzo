package ce

import "fmt"

const (
	TYPE_SYSTEM = 1
	TYPE_CUSTOM = 2
)

type Error struct {
	Msg string
}
type Err = *Error

func New(msg string) *Error {
	return &Error{
		Msg: msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}
