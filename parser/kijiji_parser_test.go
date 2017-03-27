package parser

import (
	"testing"
	"github.com/syed/kijiji-rentals/models"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/syed/kijiji-rentals/log"
	"os"
	"github.com/stretchr/testify/assert"
	"net/url"
	"fmt"
	"unicode/utf8"
)

const searchKeyword string = "Downtown"

func TestBuildQueryURL(t *testing.T) {
	query := models.KijijiQuery{Keyword: searchKeyword}

	queryURL, err := BuildQueryURL(query)

	if err != nil {
		t.Error(err)
	}

	t.Log(queryURL)

	if !strings.Contains(queryURL.String(), searchKeyword) {
		t.Error("Expected queryURL to have kewyord", searchKeyword, "Found ", queryURL)
	}
}

func TestGetKijijiPage(t *testing.T) {
	query := models.KijijiQuery{Keyword: searchKeyword}
	queryURL, err := BuildQueryURL(query)

	if err != nil {
		t.Error(err)
	}

	pageData, err := GetKijijiPage(queryURL)

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(strings.ToLower(pageData), strings.ToLower(searchKeyword)) {
		t.Error("Did not find keyword in response:", pageData)
	}
}

func TestParseListings(t *testing.T) {
	fd, err := os.Open("test_resources/b.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(fd)
	if err != nil {
		t.Error(err)
	}

	listings := ParseListings(doc)
	assert.Equal(t, len(listings), 21)
}

func TestParseAd(t *testing.T) {
	fd, err := os.Open("test_resources/c.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(fd)
	if err != nil {
		t.Error(err)
	}

	ad, err := ParseAd(doc)

	log.Debug(ad.Address)

	u, _ := url.Parse("http://www.google.com/?test=value")
	q := u.Query()
	q.Set("test", ad.Address)
	u.RawQuery = q.Encode()
	assert.Equal(t, ad.Price, 575.0)
}
