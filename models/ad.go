package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Coordinates struct {
	gorm.Model
	Lat  float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type KijijiAd struct {
        DateListed  time.Time `json:"date_listed"`
        Url         string `json:"url" gorm:"primary_key"`
        Title       string `json:"title"`
        Description string `json:"description"`
        Price       float64 `json:"price"`
        Address     string `json:"address"`
        Lat     	float64 `json:"lat"`
        Lng 		float64 `json:"lng"`
}
