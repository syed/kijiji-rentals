package main

import (
	"testing"
	"github.com/syed/kijiji-rentals/models"
	"encoding/json"
)


func TestQueryKijijiJsonParse(t *testing.T) {

	jsonStr := `{"keyword":"abcd","posted_after":"2017-04-10T04:00:00.000Z",
		"bounds":
			{"ne":{"lat":45.59116711944652,"lng":-73.28566329956055},
			"sw":{"lat":45.42491011400749,"lng":-73.82433670043946}
		}}`

	q := models.KijijiQuery{}

	if err := json.Unmarshal([]byte(jsonStr), &q); err != nil {
		t.Error(err)
	}
}
