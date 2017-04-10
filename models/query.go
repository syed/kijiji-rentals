package models

import "time"

type Point struct{
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
}

type Bounds  struct{
        Ne Point `json:"ne"`
        Sw Point `json:"nw"`
}

type KijijiQuery struct {
	MinPrice    int `json:"min_price"`
	MaxPrice    int `json:"max_price"`
	Keyword     string `json:"keyword"`
	PageNumber  int `json:"page_number"`
	PostedAfter time.Time `json:"posted_after"`
        Bounds Bounds `json:"bounds"`
}
