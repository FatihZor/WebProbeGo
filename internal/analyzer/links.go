package analyzer

import (
	"net/url"
	"strings"

	"github.com/chromedp/chromedp"
)

type RawLink struct {
	Href string `json:"href"`
	Text string `json:"text"`
	Rel  string `json:"rel"`
}

type LinkInfo struct {
	URL      string
	Text     string
	Rel      string
	Internal bool
}

// LinkActions: collect all <a href> links with href, text, rel.
func LinkActions(raw *[]RawLink) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Evaluate(`(function() {
			const anchors = Array.from(document.querySelectorAll('a[href]'));
			return anchors.map(a => ({
				href: a.href || "",
				text: (a.textContent || "").trim(),
				rel:  a.getAttribute('rel') || ""
			}));
		})()`, raw),
	}
}

// PostprocessLinks: filter, dedupe, classify internal/external.
func PostprocessLinks(origin string, raws []RawLink) ([]LinkInfo, []LinkInfo) {
	origURL, _ := url.Parse(origin)
	originHost := ""
	if origURL != nil {
		originHost = strings.ToLower(origURL.Hostname())
	}

	seen := make(map[string]struct{})
	internal := make([]LinkInfo, 0, len(raws))
	external := make([]LinkInfo, 0, len(raws))

	for _, r := range raws {
		h := strings.TrimSpace(r.Href)
		if h == "" {
			continue
		}
		// Skip non-http(s)
		lh := strings.ToLower(h)
		if strings.HasPrefix(lh, "mailto:") ||
			strings.HasPrefix(lh, "tel:") ||
			strings.HasPrefix(lh, "javascript:") ||
			strings.HasPrefix(lh, "data:") {
			continue
		}
		// Remove pure fragments
		if lh == "#" || strings.HasPrefix(lh, "#") {
			continue
		}

		u, err := url.Parse(h)
		if err != nil {
			continue
		}
		// Make absolute if needed
		if !u.IsAbs() && origURL != nil {
			u = origURL.ResolveReference(u)
		}
		if u == nil || u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
			continue
		}

		abs := u.String()
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}

		host := strings.ToLower(u.Hostname())
		itIsInternal := false
		if originHost != "" {
			itIsInternal = (host == originHost) || strings.HasSuffix(host, "."+originHost)
		}

		info := LinkInfo{
			URL:      abs,
			Text:     r.Text,
			Rel:      r.Rel,
			Internal: itIsInternal,
		}

		if itIsInternal {
			internal = append(internal, info)
		} else {
			external = append(external, info)
		}
	}

	return internal, external
}
