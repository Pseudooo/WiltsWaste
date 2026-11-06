package models

import "time"

type WasteCollectionEvent struct {
	Title string
	Date  time.Time
}

type PostcodeAddress struct {
	Address string
	Uprn    string
}
