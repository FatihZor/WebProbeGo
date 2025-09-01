package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/FatihZor/WebProbeGo/internal/config"
	"github.com/FatihZor/WebProbeGo/internal/runner"
	"github.com/FatihZor/WebProbeGo/internal/terms"
)

func main() {
	// Flags and config parsing
	cfg := config.Config{}
	flag.StringVar(&cfg.Domain, "domain", "example.com", "Target domain or URL (e.g. example.com or https://example.com)")
	flag.IntVar(&cfg.TimeoutSec, "timeout", 30, "Global timeout (seconds)")
	flag.StringVar(&cfg.FindText, "find", "", "Single keyword/phrase to search in page body")
	flag.StringVar(&cfg.FindFile, "find-file", "", "TXT file; each line is a search term")
	flag.BoolVar(&cfg.EnableNetwork, "network", false, "Enable network logging (prints request statuses)")
	flag.IntVar(&cfg.NetworkWaitSec, "network-wait", 3, "Extra wait (seconds) after body is ready (only if --network)")
	flag.StringVar(&cfg.UserAgent, "user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36",
		"Custom User-Agent string")
	flag.StringVar(&cfg.ScreenshotFile, "screenshot", "", "Save full-page screenshot to PNG (e.g. out.png)")
	flag.BoolVar(&cfg.EnableMeta, "meta", false, "Extract meta info (title/description/keywords/canonical/favicon)")

	flag.Parse()

	if cfg.Domain == "" {
		log.Fatal("Please provide a valid --domain")
	}
	if !strings.HasPrefix(cfg.Domain, "http://") && !strings.HasPrefix(cfg.Domain, "https://") {
		cfg.Domain = "https://" + cfg.Domain
	}

	// Collect search terms from --find and --find-file
	st, err := terms.Load(cfg.FindText, cfg.FindFile)
	if err != nil {
		log.Fatalf("find terms error: %v", err)
	}
	cfg.SearchTerms = st

	start := time.Now()
	if err := runner.Run(context.Background(), &cfg); err != nil {
		log.Printf("⚠️  run error: %v", err)
	}
	fmt.Printf("\n⏱️ Elapsed: %.2f seconds\n", time.Since(start).Seconds())
}
