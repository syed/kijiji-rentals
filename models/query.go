package models

import "time"

type KijijiQuery struct {
	MinPrice    int
	MaxPrice    int
	Keyword     string
	PageNumber  int
	PostedAfter time.Time
}
