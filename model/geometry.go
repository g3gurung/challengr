package model

//geometry is struct for parsing geojson coordinates for postgis
type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
