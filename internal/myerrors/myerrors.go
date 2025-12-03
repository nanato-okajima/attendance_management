package myerrors

import "fmt"

type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %v", e.Err)
}

type UnauthorizedError struct {
	Err error
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.Err)
}

type ForbiddenError struct {
	Err error
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %v", e.Err)
}

type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %v", e.Err)
}

type InternalServerError struct {
	Err error
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.Err)
}

type ConflictError struct {
	Err error
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %v", e.Err)
}
