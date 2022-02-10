package model

import (
	"errors"
)

func NewErrorJSON() (err *ErrorJSON) {
	return &ErrorJSON{}
}

type ErrorJSON struct {
	err     error
	info    string
	details string
}

func (e *ErrorJSON) Error() string {
	return e.err.Error()
}

func (e *ErrorJSON) Info() string {
	if e.info != "" {
		return e.info
	}

	return e.err.Error()
}

func (e *ErrorJSON) Details() string {
	return e.details
}

func (e *ErrorJSON) WithErrorStr(err string) (errJSON *ErrorJSON) {
	e.err = errors.New(err)

	return e
}

func (e *ErrorJSON) WithError(err error) (errJSON *ErrorJSON) {
	e.err = err

	return e
}

func (e *ErrorJSON) WithInfo(info string) (errJSON *ErrorJSON) {
	e.info = info

	return e
}

func (e *ErrorJSON) WithDetails(details string) (errJSON *ErrorJSON) {
	e.details = details

	return e
}
