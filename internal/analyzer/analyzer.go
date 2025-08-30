package analyzer

import (
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// NavigateAndGetBody navigates, waits body ready, optional delay, then returns body HTML.
func NavigateAndGetBody(url, userAgent string, extraWaitSec int) (actions []chromedp.Action, target *string) {
	var body string
	actions = append(actions,
		emulation.SetUserAgentOverride(userAgent),
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if extraWaitSec > 0 {
		actions = append(actions, chromedp.Sleep(time.Duration(extraWaitSec)*time.Second))
	}
	actions = append(actions, chromedp.InnerHTML("body", &body, chromedp.ByQuery))
	return actions, &body
}

// SearchTerms reports which terms exist in the given body.
func SearchTerms(body string, terms []string) map[string]bool {
	results := make(map[string]bool, len(terms))
	lower := strings.ToLower(body)
	for _, t := range terms {
		results[t] = strings.Contains(lower, strings.ToLower(t))
	}
	return results
}
