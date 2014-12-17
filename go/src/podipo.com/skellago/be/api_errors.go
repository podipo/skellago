package be

type APIError struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	URL     string `json:"url,omitempty"`
}

var (
	NotLoggedInError = APIError{
		Id:      "not_logged_in",
		Message: "Not logged in",
	}
	ForbiddenError = APIError{
		Id:      "forbidden",
		Message: "Forbidden for this user",
	}
	JSONParseError = APIError{
		Id:      "json_parse_error",
		Message: "JSON parse error",
	}
	IncorrectVersionError = APIError{
		Id:      "incorrect_version",
		Message: "Incorrect version",
	}
	MethodNotAllowedError = APIError{
		Id:      "method_not_allowed",
		Message: "Method not allowed",
	}
	BadRequestError = APIError{
		Id:      "bad_request",
		Message: "Bad request",
	}
	UnprocessableError = APIError{
		Id:      "unprocessable_error",
		Message: "Unprocessable",
	}
)
