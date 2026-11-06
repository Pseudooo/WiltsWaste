package models

import "time"

type WasteCollectionEvent struct {
	Title string
	Date  time.Time
}

type PostcodeAddress struct {
	Address string `json:"address"`
	Uprn    string `json:"uprn"`
}
