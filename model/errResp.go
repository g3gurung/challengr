package model

//ErrResp struct is used for send http error message
type ErrResp struct {
	Error  interface{} `json:"error"`
	Fields *[]string   `json:"fields,omitempty"`
}
