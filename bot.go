package main

import (
	"context"
	"errors"
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

var configName = flag.String("f", "config.yml", "Name of the file containing the bot configuration")

type Config struct {
	Timeout     string `yaml:"timeout"`
	PageTimeout time.Duration
	Domain      string            `yaml:"domain"`
	Cookie      map[string]string `yaml:"cookie"`
	HttpOnly    bool              `yaml:"httpOnly"`
}

func initConfig() (*Config, error) {
	c := Config{}
	if err := c.parseConfig(*configName); err != nil {
		return nil, err
	}
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, err
	}
	c.PageTimeout = timeout
	return &c, nil
}

func (c *Config) getUrl(w http.ResponseWriter, r *http.Request) {
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
	err := chromedp.Run(ctx, c.Setcookies(&res))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url[0])

	err = chromedp.Run(ctx,
		chromedp.Navigate(url[0]),
		chromedp.Sleep(c.PageTimeout),
	)
	//time.Sleep(2 * time.Second)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Config) parseConfig(fileName string) error {
	if fileName == "" {
		return errors.New("config file not found")
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return errors.New("config file not found")
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	return nil
}

func (c *Config) Setcookies(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for key, element := range c.Cookie {
				fmt.Println("Key:", key, "=>", "Element:", element)
				err := network.SetCookie(key, element).
					WithExpires(&expr).
					WithDomain(c.Domain).
					WithHTTPOnly(c.HttpOnly).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}

func main() {
	flag.Parse()
	config, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", config.getUrl)

	err = http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}
