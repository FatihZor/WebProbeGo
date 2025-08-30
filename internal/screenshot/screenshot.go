package screenshot

import (
	"os"

	"github.com/chromedp/chromedp"
)

// FullPage returns an action to capture a full-page screenshot.
func FullPage(buf *[]byte, quality int) chromedp.Action {
	return chromedp.FullScreenshot(buf, quality)
}

func SavePNG(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
