package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	http.HandleFunc("/", getUrl)

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {
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
		"session", "eyJ1c2VyIjoiYWRtaW4ifQ.ZPNAEw.Crkz32wP5psNOH1hJXvi4ePTFbw",
	))
	if err != nil {
		log.Fatal(err)
	}

	// var resbot []string
	fmt.Println(url[0])
	err = chromedp.Run(ctx,
		chromedp.Navigate(url[0]),
		// chromedp.Evaluate(`Object.keys(window);`, &resbot),
	)
	time.Sleep(2 * time.Second)

	if err != nil {
		log.Fatal(err)
	}
}

func setcookies(res *string, cookies ...string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain(os.Args[1]).
					WithHTTPOnly(true).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}
