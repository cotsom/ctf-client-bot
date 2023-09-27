// Command cookie is a chromedp example demonstrating how to set a HTTP cookie
// on requests.
package example

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx, setcookies(
		&res,
		"session", "qweqwe",
		"cookie2", "value2",
	))
	if err != nil {
		log.Fatal(err)
	}

	var resbot []string
	err = chromedp.Run(ctx,
		chromedp.Navigate(`http://localhost:8000/`),
		chromedp.Evaluate(`Object.keys(window);`, &resbot),
	)
	time.Sleep(2 * time.Second)

	if err != nil {
		log.Fatal(err)
	}
}

// setcookies returns a task to navigate to a host with the passed cookies set
// on the network request.
func setcookies(res *string, cookies ...string) chromedp.Tasks {
	if len(cookies)%2 != 0 {
		panic("length of cookies must be divisible by 2")
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain("localhost").
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
