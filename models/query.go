package models

import "time"

type Point struct{
        lat float64 `json:"lat"`
        lng float64 `json:"lng"`
}

type Bounds  struct{
        ne Point `json:"ne"`
        sw Point `json:"nw"`
}

type KijijiQuery struct {
	MinPrice    int `json:"min_price"`
	MaxPrice    int `json:"max_price"`
	Keyword     string `json:"keyword"`
	PageNumber  int `json:"page_number"`
	PostedAfter time.Time `json:"posted_after"`
        Bounds Bounds `json:"bounds"`
}
