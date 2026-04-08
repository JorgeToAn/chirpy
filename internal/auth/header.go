package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No 'Authorization' header found")
	}

	token, _ := strings.CutPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("Missing token in header")
	}
	return token, nil
}

func GetAPIKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No 'Authorization' header found")
	}

	apiKey, _ := strings.CutPrefix(authHeader, "ApiKey ")
	if apiKey == "" {
		return "", errors.New("Missing API Key in header")
	}
	return apiKey, nil
}
