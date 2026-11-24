package models

import "time"

type WasteCollectionEvent struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
}

type PostcodeAddress struct {
	Address string `json:"address"`
	Uprn    string `json:"uprn"`
}
