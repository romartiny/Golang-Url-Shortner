package response

// Response response types
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// constants of statuses
const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// OK ok response status
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error error response status
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
