package config

// Config holds configuration options for the web probe.
type Config struct {
	Domain         string
	TimeoutSec     int
	FindText       string
	FindFile       string
	SearchTerms    []string
	EnableNetwork  bool
	NetworkWaitSec int
	UserAgent      string
	ScreenshotFile string
	EnableMeta     bool
	EnableLinks    bool
}
