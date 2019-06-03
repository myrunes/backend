package webserver

type listResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}
