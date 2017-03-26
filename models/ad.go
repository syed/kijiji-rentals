package models

import (
	"time"
)

type Coordinates struct {
	Lat  float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type KijijiAd struct {
	DateListed  time.Time `json:"date_listed"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Address     string `json:"address"`
	MapLocation Coordinates `json:"map_location"`
}
