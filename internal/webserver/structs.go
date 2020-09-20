package webserver

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/objects"
)

// listResponse wraps a response of
// an arra of elements containing the
// array as Data and the length of the
// array as N.
type listResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}

// userRequest describes a request body
// for altering a user object.
type userRequest struct {
	Username        string `json:"username"`
	DisplayName     string `json:"displayname"`
	NewPassword     string `json:"newpassword"`
	CurrentPassword string `json:"currpassword"`
}

// alterFavoriteRequest describes the
// request body for modifying the array
// of favorites of a user.
type alterFavoriteRequest struct {
	Favorites []string `json:"favorites"`
}

// createShareRequest describes the
// request body for creating a page
// share.
type createShareRequest struct {
	MaxAccesses int       `json:"maxaccesses"`
	Expires     time.Time `json:"expires"`
	Page        string    `json:"page"`
}

// shareResponse wraps the response
// data when requesting a page share
// object containing the data for the
// share and the liquified data for
// the page which is shared and the
// user which owns the page.
type shareResponse struct {
	Share *objects.SharePage `json:"share"`
	Page  *objects.Page      `json:"page"`
	User  *objects.User      `json:"user"`
}

// pageOrderRequest describes the request
// when modifying the users page order.
type pageOrderRequest struct {
	PageOrder []snowflake.ID `json:"pageorder"`
}

// setMailRequest describes the reuqest
// model for setting or resetting a
// users e-mail specification.
type setMailRequest struct {
	MailAddress     string `json:"mailaddress"`
	Reset           bool   `json:"reset"`
	CurrentPassword string `json:"currpassword"`
}

// confirmMail desribes the request model
// to confirm a mail settings change.
type confirmMail struct {
	Token string `json:"token"`
}

// passwordReset describes the request model
// on requesting a password reset mail.
type passwordReset struct {
	MailAddress string `json:"mailaddress"`
}

// confirmPasswordReset describes the password
// reset model containing the generated
// verification, the new password string and
// the page name guesses.
type confirmPasswordReset struct {
	reCaptchaResponse

	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// mailConfirmtationData wraps the mail address
// and user ID which is saved in the mail
// confirmation cache to check and identify
// e-mail confirmations.
type mailConfirmationData struct {
	UserID      snowflake.ID
	MailAddress string
}

// reCaptchaResponse wraps a ReCAPTCHA response
// token for ReCAPTCHA validation.
type reCaptchaResponse struct {
	ReCaptchaResponse string `json:"recaptcharesponse"`
}
