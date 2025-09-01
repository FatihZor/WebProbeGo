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

	// 2) Meta (opsiyonel)
	var meta analyzer.MetaInfo
	if cfg.EnableMeta {
		tasks = append(tasks, analyzer.MetaActions(&meta)...)
	}

	// 3) Screenshot (opsiyonel)
	var shotBuf []byte
	if cfg.ScreenshotFile != "" {
		tasks = append(tasks, screenshot.FullPage(&shotBuf, 90))
	}

	// Run actions
	if err := chromedp.Run(ctx, tasks...); err != nil {
		// Fallback sadece body
		fallbackCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
		defer cancel2()
		_ = chromedp.Run(fallbackCtx, chromedp.InnerHTML("body", bodyPtr, chromedp.ByQuery))
	}

	// Arama sonuçları (varsa)
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

	// Meta sonuçlarını yaz
	if cfg.EnableMeta {
		meta.Clean()
		fmt.Println("\n🧭 Meta Info")
		fmt.Println("Title      :", meta.Title)
		fmt.Println("Description:", meta.Description)
		fmt.Println("Keywords   :", meta.Keywords)
		fmt.Println("Canonical  :", meta.Canonical)
		fmt.Println("Favicon    :", meta.FaviconURL)
	}

	// Screenshot kaydet
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
