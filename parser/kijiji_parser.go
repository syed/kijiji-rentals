package parser

import (
	"net/url"
	"github.com/syed/kijiji-rentals/models"
	"github.com/syed/kijiji-rentals/log"
	"errors"
	"strconv"
	"net/http"
	"net/http/cookiejar"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
	"math/rand"
	"encoding/json"
	"sync"
)

const (
	BASE_URL            string = "http://www.kijiji.ca"
	BASE_SEARCH_URL     string = "http://www.kijiji.ca/b-search.html?formSubmit=true&ll=&categoryId=0&categoryName=appartements%2C+condos&locationId=1700281&pageNumber=1&minPrice=&maxPrice=&adIdRemoved=&sortByName=dateDesc&userId=&origin=&searchView=LIST&urgentOnly=false&cpoOnly=false&carproofOnly=false&highlightOnly=false&gpTopAd=false&adPriceType=&brand=&keywords=charlevoix&SearchCategory=0&SearchLocationPicker=Ville+de+Montr%C3%A9al&siteLocale=en_CA"
	MAX_CONCURRENT_REQS int    = 5
	REQ_TIMEOUT                = 2 * time.Minute
	MAX_PAGES           int    = 5
)

var jar *cookiejar.Jar
var err error

func init()  {
	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Warning("Unable to create cookie jar")
		panic(err)
	}
}

/* Given a keyword, searches kijiji and returns the raw HTML */
func BuildQueryURL(query models.KijijiQuery) (*url.URL, error) {
	if len(query.Keyword) == 0 {
		return nil, errors.New("Keyword must be specified")
	}

	searchURL, _ := url.ParseRequestURI(BASE_SEARCH_URL)
	urlQuery := searchURL.Query()

	urlQuery.Set("keywords", query.Keyword)

	urlQuery.Set("pageNumber", strconv.Itoa(query.PageNumber))

	if query.MinPrice > 0 {
		urlQuery.Set("minPrice", strconv.Itoa(query.MinPrice))
	}

	if query.MinPrice > 0 {
		urlQuery.Set("maxPrice", strconv.Itoa(query.MaxPrice))
	}

	searchURL.RawQuery = urlQuery.Encode()
	return searchURL, nil
}

func GetKijijiPage(url *url.URL) (string, error) {
	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Get(url.String())

	if err != nil {
		log.Warning("Unable to get response for search query", url)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SearchKijiji(query models.KijijiQuery) ([]models.KijijiAd, error) {

	ads := make([]models.KijijiAd, 0)


	for i := 1; i <= MAX_PAGES; i++ {
		query.PageNumber = i
		url, err := BuildQueryURL(query)
		if err != nil {
			return nil, err
		}

		data, err := GetKijijiPage(url)
		if err != nil {
			log.Warning("Unable to get Data for URL", url)
			continue
		}

		newAds, err := ParseKijjiPage(data)
		if err != nil {
			log.Warning("Unable to parse ads for URL", url)
		}

		ads = append(ads, newAds...)

		//check if the last ad is beyond the end date
		lastAd := ads[len(ads)-1]
		if lastAd.DateListed.Before(query.PostedAfter) {
			break
		}
	}

	ads = FetchAddress(ads)
	return ads, nil
}

func ParseKijjiPage(pageHtml string) ([]models.KijijiAd, error) {

	//1. get all the links
	//2. open them in seperate goroutines
	//3. wait for their response
	//4. parse each of them and return them in an array

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHtml))

	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}

	listings := ParseListings(doc)

	result := make(chan models.KijijiAd)
	start := make(chan bool, MAX_CONCURRENT_REQS)
	ads := make([]models.KijijiAd, 0)

	for _, v := range listings {
		start <- true
		go FetchAndParseKijijiAd(v, start, result)
	}

	for i := 0; i < len(listings); i++ {
		select {
		case r := <-result:
			ads = append(ads, r)
		case <-time.After(REQ_TIMEOUT):
			break
		}
	}

	return ads, nil
}

func FetchAndParseKijijiAd(l string, start chan bool, res chan models.KijijiAd) {

	u, err := url.ParseRequestURI(l)
	resp, err := GetKijijiPage(u)
	time.AfterFunc(2*time.Second, func() {
		<-start
	})

	if err != nil {
		log.Warning("Error getting URL:", u)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))

	if err != nil {
		log.Warning(err.Error())
		return
	}

	ad, err := ParseAd(doc)
	ad.Url = u.String()

	res <- ad
}

func ParseListings(doc *goquery.Document) []string {

	listings := make([]string, 0)

	doc.Find(".info-container").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a").Attr("href")
		if exists {
			listings = append(listings, BASE_URL+link)
		}
	})

	for i, v := range listings {
		log.Debug(i, v)
	}

	return listings
}

func ParseAd(doc *goquery.Document) (models.KijijiAd, error) {
	ad := models.KijijiAd{}

	ad.Title = doc.Find("h1").Text()
	ad.Description = doc.Find("#UserContent").Text()

	doc.Find(".ad-attributes tr").Each(func(i int, s *goquery.Selection) {
		attr := strings.ToLower(s.Find("th").Text())
		val := s.Find("td").Text()

		if attr == "price" {
			val = strings.Trim(val, " $\n")
			val = strings.Replace(val, ",", "", -1)
			price, err := strconv.ParseFloat(val, 64)
			if err == nil {
				ad.Price = price
			} else {
				log.Debug("Invalid price ", val)
			}
		} else if attr == "address" {
			val := strings.Replace(val, "View map", "", 1)
			ad.Address = strings.Trim(val, " \n")
		} else if attr == "date listed" {
			layout := "02-Jan-06"
			t, err := time.Parse(layout, val)
			if err == nil {
				ad.DateListed = t
			} else {
				log.Debug("Invalid date ", val)
			}
		}
	})

	return ad, nil
}

func FetchAddress(ads []models.KijijiAd) []models.KijijiAd {
	mapUrl := "http://maps.google.com/maps/api/geocode/json?address=ReplaceAddress"

	var wg sync.WaitGroup
	wg.Add(len(ads))

	for i := 0; i < len(ads); i++ {
		go func(i int, ads []models.KijijiAd) {
			defer wg.Done()
			addressUrl, err := url.Parse(mapUrl)
			if err != nil {
				log.Warning(err)
				return
			}

			query := addressUrl.Query()
			query.Set("address", ads[i].Address)
			addressUrl.RawQuery = query.Encode()

			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			resp, err := http.Get(addressUrl.String())
			if err != nil {
				log.Warning(err)
				return
			}
			defer resp.Body.Close()

			var location models.LocationResult
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&location)
			if err != nil {
				log.Warning(err)
				return
			}

			if len(location.Results) == 0 {
				log.Warning("No results for address:", addressUrl)
				return
			}

			ads[i].MapLocation.Lat = location.Results[0].Geometry.Location.Lat
			ads[i].MapLocation.Lng = location.Results[0].Geometry.Location.Lng

		}(i, ads)
	}
	wg.Wait()
	return ads
}
