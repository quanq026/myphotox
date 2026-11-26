package hub

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/photoprism/photoprism/pkg/http/safe"
)

// Default service base URLs disabled in MyPhotoX standalone mode.
const (
	ProdBaseURL = ""
	TestBaseURL = ""
)

// baseURL specifies the service endpoint URL.
var baseURL = ""

// GetServiceURL returns the currently configured Hub endpoint, optionally
// appending the provided API key. An empty string is returned when Hub
// requests are disabled.
func GetServiceURL(key string) string {
	if baseURL == "" {
		return ""
	}

	if key == "" {
		return baseURL
	}

	return fmt.Sprintf(baseURL+"/%s", key)
}

// GetFeedbackServiceURL builds the feedback endpoint corresponding to the
// supplied API key. A disabled Hub service results in an empty string.
func GetFeedbackServiceURL(key string) string {
	if key == "" {
		return ""
	}

	u := GetServiceURL(key)

	if u == "" {
		return ""
	}

	return u + "/feedback"
}

// GetServiceHost extracts the Hub host name from the active base URL, or
// returns an empty string when the service is disabled or invalid.
func GetServiceHost() string {
	s := GetServiceURL("")

	if s == "" {
		return ""
	}

	u, err := url.Parse(s)

	if err != nil {
		log.Warn(err)
		return ""
	}

	return u.Host
}

// requestError returns the error of a failed service request with the service host and the underlying cause
// instead of the request URL.
func requestError(err error) error {
	var urlErr *url.Error

	if !errors.As(err, &urlErr) || urlErr == nil {
		return err
	}

	desc := "service request failed"

	if host := GetServiceHost(); host != "" {
		desc = fmt.Sprintf("request to %s failed", host)
	}

	if urlErr.Err == nil {
		return errors.New(desc)
	}

	return fmt.Errorf("%s (%w)", desc, urlErr.Err)
}

// ValidateServiceURL ensures outbound Hub requests use an absolute HTTPS URL.
func ValidateServiceURL(rawURL string) error {
	u, err := safe.URL(rawURL)
	if err != nil {
		return err
	}

	if strings.ToLower(u.Scheme) != "https" {
		return fmt.Errorf("unsupported service URL scheme")
	}

	return nil
}

// SetBaseURL is a no-op in MyPhotoX standalone mode so outbound Hub calls remain disabled.
func SetBaseURL(_ string) {}

// Disabled reports whether outbound Hub requests have been switched off.
func Disabled() bool {
	return baseURL == ""
}

// Disable clears the Hub endpoint so no network calls are attempted.
func Disable() {
	// Return if already disabled.
	if Disabled() {
		return
	}

	// Remove configured endpoint URL to disable service.
	baseURL = ""
	log.Debugf("config: disabled hub service requests")
}

// ApplyTestConfig disables outbound Hub service requests in MyPhotoX standalone mode.
func ApplyTestConfig() {
	Disable()
}
