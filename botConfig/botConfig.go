package botConfig

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Domain   string            `yaml:"domain"`
	Cookie   map[string]string `yaml:"cookie"`
	HttpOnly bool              `yaml:"httpOnly"`
}

func setcookies(res *string, domain string, cookies map[string]string, httpOnly bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for key, element := range cookies {
				fmt.Println("Key:", key, "=>", "Element:", element)
				err := network.SetCookie(key, element).
					WithExpires(&expr).
					WithDomain(domain).
					WithHTTPOnly(httpOnly).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}

func parseyaml(fileName string) (string, map[string]string, bool) {
	if fileName == "" {
		fmt.Println("config file not found")
		// return
	}

	yamlFile, err := ioutil.ReadFile(fileName)
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
	httpOnly := yamlConfig.HttpOnly

	return domain, cookie, httpOnly
}
