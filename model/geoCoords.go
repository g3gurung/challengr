package model

type geoCoords struct {
	ID   int64   `json:"id" sql:"id"`
	Long float64 `json:"long" sql:"long"`
	Lat  float64 `json:"lat" sql:"lat"`
}
