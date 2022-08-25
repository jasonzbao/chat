package dynaerrors

import "github.com/pkg/errors"

type DynaError struct {
	Err error `json:"err"`
}

func (e DynaError) Error() string {
	return e.Err.Error()
}

func (e DynaError) Unwrap() error {
	return e.Err
}

var ErrorNameNotSet = &DynaError{
	Err: errors.New("Name must be set before messages are sent"),
}

var ErrorInvalidInstruction = &DynaError{
	Err: errors.New("Instruction is not recognized"),
}

var ErrorExitChat = &DynaError{
	Err: errors.New("Please exit chat"),
}
