package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/syed/kijiji-rentals/models"
        "fmt"
)

func SaveAdsToDB(ads []models.KijijiAd) {

	db, err := gorm.Open("sqlite3", "kijiji.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.KijijiAd{})

	dbRecord := models.KijijiAd{}
	for i := range ads {
		ad := ads[i]

		db.Where("url = ?", ad.Url).First(&dbRecord)

		if dbRecord.Url != ad.Url {
			db.Create(ads[i])
		}
	}
}

func GetAdsFromDB(query models.KijijiQuery) []models.KijijiAd {

	db, err := gorm.Open("sqlite3", "kijiji.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

        ads := make([]models.KijijiAd, 0, 0)
        db.Where("date_listed >= ? AND description LIKE ?",
                query.PostedAfter,
                fmt.Sprintf("%%%s%%", query.Keyword)).Find(&ads)

        return ads
}
