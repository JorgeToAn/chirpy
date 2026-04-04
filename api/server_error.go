package api

type ServerError int

const (
	ErrorGeneric ServerError = iota
	ErrorUnauthorized
	ErrorBadCredentials
)

var errorMessage = map[ServerError]string{
	ErrorGeneric:        "Something went wrong",
	ErrorUnauthorized:   "Unauthorized",
	ErrorBadCredentials: "Incorrect email or password",
}

func (se ServerError) String() string {
	return errorMessage[se]
}
