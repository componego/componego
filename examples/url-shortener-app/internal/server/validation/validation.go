package validation

import (
	"net/url"
)

func IsValidUrl(value string) bool {
	if _, err := url.ParseRequestURI(value); err != nil {
		return false
	}
	parsedUrl, err := url.Parse(value)
	if err != nil || parsedUrl.Scheme == "" {
		return false
	}
	return true
}
