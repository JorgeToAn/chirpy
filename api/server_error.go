package api

type ServerError int

const (
	ErrorGeneric ServerError = iota
	ErrorUnauthorized
	ErrorBadCredentials
	ErrorForbidden
)

var errorMessage = map[ServerError]string{
	ErrorGeneric:        "Something went wrong",
	ErrorUnauthorized:   "Unauthorized",
	ErrorBadCredentials: "Incorrect email or password",
	ErrorForbidden:      "Forbidden",
}

func (se ServerError) String() string {
	return errorMessage[se]
}
