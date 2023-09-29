package main

import (
	"bot/config"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var configName = flag.String("f", "config.yml", "Name of the file containing the bot configuration")

func main() {
	flag.Parse()

	http.HandleFunc("/", getUrl)

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {

	configFile, err := config.Parseyaml(*configName)

	if err != nil {
		log.Fatal(err)
	}

	query := r.URL.Query()
	url, present := query["url"]
	if !present || len(url) == 0 {
		fmt.Println("url not present")
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(strings.Join(url, ",")))
	if !strings.Contains(url[0], "://") {
		url[0] = "http://" + url[0]
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err = chromedp.Run(ctx, config.Setcookies(
		&res,
		*configFile,
	))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url[0])
	err = chromedp.Run(ctx,
		chromedp.Navigate(url[0]),
	)
	time.Sleep(2 * time.Second)

	if err != nil {
		log.Fatal(err)
	}
}
