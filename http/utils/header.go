package utils

import (
	"fmt"
	"net/http"
	"regexp"
)

var ssoIDDigitPattern = regexp.MustCompile(`^(svc)?[0-9]{9}$`)

// GetSsoID extracts the current SSO ID from in coming request.
func GetSsoID(req *http.Request) (string, error) {
	ssoRaw := req.Header.Get("Sso")
	if !ssoIDDigitPattern.MatchString(ssoRaw) {
		return "", fmt.Errorf("invalid SSO value")
	}
	return ssoRaw, nil
}
