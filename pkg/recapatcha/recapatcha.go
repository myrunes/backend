// package recapatcha provides bindings to google's
// ReCAPTCHA validation service.
package recapatcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	endpoint    = "https://www.google.com/recaptcha/api/siteverify"
	contentType = "application/json"
)

// Response wraps the HTTP response from the
// validation endpoint containing success state,
// challenge timestamp, hostname and possible
// error codes on validation failure.
//
// Read here for more information:
// https://developers.google.com/recaptcha/docs/verify
type Response struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// Validate sends a HTTP POST request to google's ReCAPTCHA
// validation endpoint using the specified secret, response
// and remoteIP (optional) as parameters.
//
// The prased Response result object is returned. Errors are
// only returned when the response parsing or the request
// itself failed. Invalid validation will not return an error.
func Validate(secret, response string, remoteIP ...string) (res *Response, err error) {
	remoteAddrParam := ""
	if len(remoteIP) > 0 {
		remoteAddrParam = "&remoteip=" + remoteIP[0]
	}

	url := fmt.Sprintf("%s?secret=%s&response=%s%s", endpoint, secret, response, remoteAddrParam)

	var httpRes *http.Response
	if httpRes, err = http.Post(url, contentType, nil); err != nil {
		return
	}

	if httpRes.StatusCode != 200 {
		err = fmt.Errorf("response code was %d", httpRes.StatusCode)
		return
	}

	res = new(Response)
	err = json.NewDecoder(httpRes.Body).Decode(res)

	return
}
