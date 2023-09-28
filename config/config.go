package config

import (
	"context"
	"errors"
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

func Setcookies(res *string, config YamlConfig) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for key, element := range config.Cookie {
				fmt.Println("Key:", key, "=>", "Element:", element)
				err := network.SetCookie(key, element).
					WithExpires(&expr).
					WithDomain(config.Domain).
					WithHTTPOnly(config.HttpOnly).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}

func Parseyaml(fileName string) (*YamlConfig, error) {
	if fileName == "" {
		return nil, errors.New("config file not found")
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, errors.New("config file not found")
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	// domain := yamlConfig.Domain
	// cookie := yamlConfig.Cookie
	// httpOnly := yamlConfig.HttpOnly

	return &yamlConfig, nil
}
