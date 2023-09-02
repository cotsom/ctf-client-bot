package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Domain string            `yaml:"domain"`
	Cookie map[string]string `yaml:"cookie"`
}

func main() {
	http.HandleFunc("/", getUrl)

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	domain, cookie := parseyaml()
	fmt.Println(domain, cookie)

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

func setcookies(res *string, domain string, cookies ...string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain(domain).
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

func parseyaml() (string, map[string]string) {
	var fileName = flag.String("f", "", "Значение флага -f")
	flag.Parse()

	if *fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")
		// return
	}

	yamlFile, err := ioutil.ReadFile(*fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		// return
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	domain := yamlConfig.Domain
	cookie := yamlConfig.Cookie

	return domain, cookie
}
