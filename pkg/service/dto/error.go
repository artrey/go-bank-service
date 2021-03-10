package dto

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MakeError(code string, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

func MakeUnknownError(err error) Error {
	return MakeError("unknown", err)
}
