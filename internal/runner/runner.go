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

func Run(parent context.Context, cfg *config.Config) error {
	ctx, cancel := chromedp.NewContext(parent)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.TimeoutSec)*time.Second)
	defer cancel()

	var tasks []chromedp.Action

	if cfg.EnableNetwork {
		network.AttachLogger(ctx)
		tasks = append(tasks, network.EnableNetwork())
	}

	// 1) Navigate & Body
	bodyActions, bodyPtr := analyzer.NavigateAndGetBody(cfg.Domain, cfg.UserAgent, ifThen(cfg.EnableNetwork, cfg.NetworkWaitSec))
	tasks = append(tasks, bodyActions...)

	// 2) Meta (optional)
	var meta analyzer.MetaInfo
	if cfg.EnableMeta {
		tasks = append(tasks, analyzer.MetaActions(&meta)...)
	}

	// 3) Links (optional)
	var rawLinks []analyzer.RawLink
	if cfg.EnableLinks {
		tasks = append(tasks, analyzer.LinkActions(&rawLinks)...)
	}

	// 4) Screenshot (optional)
	var shotBuf []byte
	if cfg.ScreenshotFile != "" {
		tasks = append(tasks, screenshot.FullPage(&shotBuf, 90))
	}

	// Run actions
	if err := chromedp.Run(ctx, tasks...); err != nil {
		// fallback: body
		fallbackCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
		defer cancel2()
		_ = chromedp.Run(fallbackCtx, chromedp.InnerHTML("body", bodyPtr, chromedp.ByQuery))
	}

	// Search results
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

	// Write meta info
	if cfg.EnableMeta {
		meta.Clean()
		fmt.Println("\n🧭 Meta Info")
		fmt.Println("Title      :", meta.Title)
		fmt.Println("Description:", meta.Description)
		fmt.Println("Keywords   :", meta.Keywords)
		fmt.Println("Canonical  :", meta.Canonical)
		fmt.Println("Favicon    :", meta.FaviconURL)
	}

	// Write link inventory
	if cfg.EnableLinks {
		internals, externals := analyzer.PostprocessLinks(cfg.Domain, rawLinks)
		fmt.Printf("\n🔗 Link Inventory\n")
		fmt.Printf("Internal: %d\n", len(internals))
		fmt.Printf("External: %d\n", len(externals))

		// TODO: all links to file?
		// Show first 10 internal and first 10 external

		max := func(n, lim int) int {
			if n < lim {
				return n
			}
			return lim
		}
		fmt.Println("\nFirst 10 internal:")
		for i := 0; i < max(len(internals), 10); i++ {
			fmt.Printf("  - %s\n", internals[i].URL)
		}
		fmt.Println("\nFirst 10 external:")
		for i := 0; i < max(len(externals), 10); i++ {
			fmt.Printf("  - %s\n", externals[i].URL)
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

func ifThen(cond bool, v int) int {
	if cond {
		return v
	}
	return 0
}
