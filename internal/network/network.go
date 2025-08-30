package network

import (
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// EnableNetwork returns actions to enable network domain.
func EnableNetwork() chromedp.Action {
	return network.Enable()
}

// AttachLogger wires a listener that prints status codes.
func AttachLogger(ctx chromedp.Context) {
	go func() {
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			if resp, ok := ev.(*network.EventResponseReceived); ok {
				url := resp.Response.URL
				status := resp.Response.Status
				if status != 200 {
					fmt.Printf("⚠️  %d %s\n", status, url)
				} else {
					fmt.Printf("✅  %d %s\n", status, url)
				}
			}
		})
	}()
}
