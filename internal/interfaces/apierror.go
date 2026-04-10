package interfaces

type APIError struct {
	Code          int
	Msg           string
	OriginalError error
}

func (e *APIError) Error() string {
	return e.Msg
}
