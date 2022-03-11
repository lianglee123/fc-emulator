package utils

import (
	"errors"
	"fmt"
)

func NewError(errInfos ...interface{}) error {
	return errors.New(fmt.Sprintln(errInfos...))
}
