package models

type GetCarData struct {
	Data []CarDatum `json:"data"`
}

type CarDatum struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
