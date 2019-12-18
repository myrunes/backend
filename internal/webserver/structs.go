package webserver

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/myrunes/internal/objects"
)

type listResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}

type userRequest struct {
	Username        string `json:"username"`
	DisplayName     string `json:"displayname"`
	NewPassword     string `json:"newpassword"`
	CurrentPassword string `json:"currpassword"`
}

type alterFavoriteRequest struct {
	Favorites []string `json:"favorites"`
}

type createShareRequest struct {
	MaxAccesses int       `json:"maxaccesses"`
	Expires     time.Time `json:"expires"`
	Page        string    `json:"page"`
}

type shareResponse struct {
	Share *objects.SharePage `json:"share"`
	Page  *objects.Page      `json:"page"`
	User  *objects.User      `json:"user"`
}

type sessionsResponse struct {
	listResponse

	CurrentlyConnectedID string `json:"currentlyconnectedid"`
}

type pageOrderRequest struct {
	PageOrder []snowflake.ID `json:"pageorder"`
}

type setMailRequest struct {
	MailAddress string `json:"mailaddress"`
	Reset       bool   `json:"reset"`
}

type confirmMail struct {
	Token string `json:"token"`
}

type passwordReset struct {
	MailAddress string `json:"mailaddress"`
}

type confirmPasswordReset struct {
	Token       string   `json:"token"`
	NewPassword string   `json:"new_password"`
	PageNames   []string `json:"page_names"`
}
