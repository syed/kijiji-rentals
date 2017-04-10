package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/syed/kijiji-rentals/models"
	"github.com/syed/kijiji-rentals/log"
        "github.com/syed/kijiji-rentals/scraper"
        "github.com/syed/kijiji-rentals/db"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	index, err := ioutil.ReadFile("./static/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(index)
}

func queryKijiji(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Form", r.Form)
	fmt.Println("Got Query ")

        decoder := json.NewDecoder(r.Body)
        query  := models.KijijiQuery{}

        err := decoder.Decode(&query)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

        defer r.Body.Close()

        ads := db.GetAdsFromDB(query)

	log.Debug(fmt.Sprintf("Got %d ads ", len(ads)))

	v, err := json.Marshal(ads)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(v) // send data to client side
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/query", queryKijiji)
	http.HandleFunc("/static/", serveStatic)
	http.HandleFunc("/", serveIndex)

	log.Debug("Serving ...")

        //start scraper

        go scraper.StartScrape()


	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Error("ListenAndServe: ", err)
	}
}
