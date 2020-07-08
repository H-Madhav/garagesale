package web

// ErrorResponse return clinet error
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is use to add more info to web errors
type Error struct {
	Err    error
	Status int
}

//NewRequestError used for know error
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
