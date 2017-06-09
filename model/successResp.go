package model

//SuccessResp struct is used for sending http success response to the client
type SuccessResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
