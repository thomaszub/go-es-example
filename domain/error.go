package domain

import "fmt"

type DomainError struct {
	Reason string
}

func (e *DomainError) Error() string {
	return e.Reason
}

func NewDomainError(format string, a ...any) *DomainError {
	return &DomainError{Reason: fmt.Sprintf(format, a...)}
}

type AccountNotFoundError struct {
	Reason string
}

func (e *AccountNotFoundError) Error() string {
	return e.Reason
}

func NewAccountNotFoundError(format string, a ...any) *AccountNotFoundError {
	return &AccountNotFoundError{Reason: fmt.Sprintf(format, a...)}
}
