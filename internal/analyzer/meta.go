package analyzer

import (
	"strings"

	"github.com/chromedp/chromedp"
)

type MetaInfo struct {
	Title       string
	Description string
	Keywords    string
	Canonical   string
	FaviconURL  string
}

// MetaActions hazır DOM üzerinden gerekli alanları toplar.
// WaitReady('body') sonrası çağrılacak şekilde sadece Evaluate/Value alır.
func MetaActions(mi *MetaInfo) []chromedp.Action {
	return []chromedp.Action{
		// Title
		chromedp.Evaluate(`document.title || ""`, &mi.Title),

		// Description (priority: meta[name=description] -> og:description)
		chromedp.Evaluate(`(function(){
			const m = document.querySelector("meta[name='description']")?.content || "";
			if (m && m.trim()) return m.trim();
			const og = document.querySelector("meta[property='og:description']")?.content || "";
			return (og||"").trim();
		})()`, &mi.Description),

		// Keywords
		chromedp.Evaluate(`(document.querySelector("meta[name='keywords']")?.content || "").trim()`, &mi.Keywords),

		// Canonical
		chromedp.Evaluate(`(document.querySelector("link[rel='canonical']")?.href || "").trim()`, &mi.Canonical),

		// Favicon (priority: rel~='icon' -> /favicon.ico)
		chromedp.Evaluate(`(function(){
			// Try common rels in order
			const sels = ["link[rel='icon']","link[rel='shortcut icon']","link[rel='apple-touch-icon']","link[rel*='icon']"];
			for (const s of sels){
				const el = document.querySelector(s);
				if (el && el.href) return el.href;
			}
			// Fallback
			try {
				const u = new URL(location.href);
				return u.origin + "/favicon.ico";
			} catch(e) { return "/favicon.ico"; }
		})()`, &mi.FaviconURL),
	}
}

// Clean trims and normalizes some fields.
func (m *MetaInfo) Clean() {
	m.Title = strings.TrimSpace(m.Title)
	m.Description = strings.TrimSpace(m.Description)
	m.Keywords = strings.TrimSpace(m.Keywords)
	m.Canonical = strings.TrimSpace(m.Canonical)
	m.FaviconURL = strings.TrimSpace(m.FaviconURL)
}
