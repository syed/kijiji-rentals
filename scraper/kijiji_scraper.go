package scraper

import (
	"github.com/syed/kijiji-rentals/parser"
	"github.com/syed/kijiji-rentals/models"
	"time"
        "github.com/syed/kijiji-rentals/db"
        "github.com/syed/kijiji-rentals/log"
)

const POLL_TIME time.Duration = 1 * time.Hour // Poll kijiji every hour

func StartScrape() {

	// fetch the kijiji apartment listing page
	// get listing of ads and the pagination
	// parse each ad and save it in the DB

        for {

                <-time.After(POLL_TIME)

                yesterday := time.Now().AddDate(0, 0, -1);
                query := models.KijijiQuery{Keyword: "", PostedAfter: yesterday}
                ads, err := parser.SearchKijiji(query)

                if err != nil {
                        log.Warning("Unable to fetch ads from kijiji", err.Error())
                }

                db.SaveAdsToDB(ads)

        }
}