package myerrors

type AppError interface {
	error
}

type BadRequestError struct {
	Err error
}

type InternalServerError struct {
	Err error
}

func (br *BadRequestError) Error() string {
	return br.Err.Error()
}

func (br *InternalServerError) Error() string {
	return br.Err.Error()
}
