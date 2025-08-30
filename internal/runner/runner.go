package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/FatihZor/WebProbeGo/internal/analyzer"
	"github.com/FatihZor/WebProbeGo/internal/config"
	"github.com/FatihZor/WebProbeGo/internal/network"
	"github.com/FatihZor/WebProbeGo/internal/screenshot"
	"github.com/chromedp/chromedp"
)

// Run wires up chromedp and orchestrates the flow.
func Run(parent context.Context, cfg *config.Config) error {
	// Create Chrome context
	ctx, cancel := chromedp.NewContext(parent)
	defer cancel()

	// Global timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.TimeoutSec)*time.Second)
	defer cancel()

	// Network logging (optional)
	var tasks []chromedp.Action
	if cfg.EnableNetwork {
		network.AttachLogger(ctx)
		tasks = append(tasks, network.EnableNetwork())
	}

	// Navigate & fetch body
	bodyActions, bodyPtr := analyzer.NavigateAndGetBody(cfg.Domain, cfg.UserAgent,
		ifThen(cfg.EnableNetwork, cfg.NetworkWaitSec))
	tasks = append(tasks, bodyActions...)

	// Screenshot (optional)
	var shotBuf []byte
	if cfg.ScreenshotFile != "" {
		tasks = append(tasks, screenshot.FullPage(&shotBuf, 90))
	}

	// Run
	if err := chromedp.Run(ctx, tasks...); err != nil {
		// fallback: try only body fetch with short window
		fallbackCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
		defer cancel2()
		_ = chromedp.Run(fallbackCtx, chromedp.InnerHTML("body", bodyPtr, chromedp.ByQuery))
	}

	// Report search results
	if len(cfg.SearchTerms) > 0 {
		results := analyzer.SearchTerms(*bodyPtr, cfg.SearchTerms)
		for term, ok := range results {
			if ok {
				fmt.Printf("🔍 '%s' found\n", term)
			} else {
				fmt.Printf("❌ '%s' not found\n", term)
			}
		}
	}

	// Save screenshot
	if cfg.ScreenshotFile != "" && len(shotBuf) > 0 {
		if err := screenshot.SavePNG(cfg.ScreenshotFile, shotBuf); err != nil {
			return fmt.Errorf("save screenshot: %w", err)
		}
		fmt.Printf("📸 Screenshot saved: %s\n", cfg.ScreenshotFile)
	}

	return nil
}

// ifThen returns v if cond is true, otherwise false.
func ifThen(cond bool, v int) int {
	if cond {
		return v
	}
	return 0
}
