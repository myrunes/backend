package webserver

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
