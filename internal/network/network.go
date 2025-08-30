package network

import (
	"context" // <-- EDIT: added this import for context.Context
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func EnableNetwork() chromedp.Action {
	return network.Enable()
}

// AttachLogger attaches a network event listener to log HTTP responses.
func AttachLogger(ctx context.Context) {
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
