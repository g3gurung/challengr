package model

//OneSignal struct is a model/schema for one_signal table
type OneSignal struct {
	ID       int64  `json:"id" sql:"id"`
	UserID   int64  `json:"user_id" sql:"id"`
	IMEI     string `json:"imei" sql:"imei"`
	PlayerID string `json:"player_id" sql:"playerID"`
}
