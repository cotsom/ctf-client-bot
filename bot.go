package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type YamlConfig struct {
	Domain   string            `yaml:"domain"`
	Cookie   map[string]string `yaml:"cookie"`
	HttpOnly bool              `yaml:"httpOnly"`
}

func main() {
	if len(os.Args) < 1 {

	}

	http.HandleFunc("/", getUrl)

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	configName := os.Args[1]

	domain, cookie, httpOnly := botConfig.parseyaml(configName)
	fmt.Println(domain, cookie, httpOnly)

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

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx, setcookies(
		&res,
		domain,
		cookie,
		httpOnly,
	))
	if err != nil {
		log.Fatal(err)
	}

	// var resbot []string
	fmt.Println(url[0])
	err = chromedp.Run(ctx,
		chromedp.Navigate(url[0]),
	)
	time.Sleep(2 * time.Second)

	if err != nil {
		log.Fatal(err)
	}
}
