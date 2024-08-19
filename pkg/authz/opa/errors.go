package opa

import (
	"fmt"
)

type ErrOpaFailure struct {
	Code     string
	Row      int
	Message  string
	Filename string
}

func (e ErrOpaFailure) Error() string {
	return fmt.Sprintf("%s:%d - %s", e.Filename, e.Row, e.Message)
}

type ErrOpaFailures []ErrOpaFailure

func (e ErrOpaFailures) Error() string {
	return fmt.Sprintf("found %d opa failures", len(e))
}
