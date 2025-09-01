# WebProbeGo

**WebProbeGo** — a fast and versatile web analysis tool built with Go and Chromedp.  
It allows you to probe websites, capture screenshots, inspect network requests, and search for keywords — all from a simple CLI.

---

## ✨ Features

- **Network request logging** (`--network`)  
  Capture all requests a page makes and report their HTTP status codes.

- **Keyword/phrase search** (`--find`, `--find-file`)  
  Search for a single term or a list of terms (from a `.txt` file) in the page body.

- **Meta info extraction** (`--meta`)  
  Extract title, description, keywords, canonical URL, and favicon.

- **Link inventory** (`--links`)  
  Collect all `<a>` links and classify them as internal vs external.

- **Screenshots** (`--screenshot`)  
  Save a full-page PNG screenshot of the target site.

- **Execution time measurement**  
  Report how long the operation took in seconds.

---

## 🚧 Roadmap

Planned improvements:

- Heading map (H1–H3, optionally deeper)
- Broken link checker (`<a>` and `<img>`)
- SSL/expiry check
- Redirect chain report
- Selector-only screenshots (e.g. `--screenshot="#main"`)
- Mobile/tablet viewport emulation
- Security headers report (CSP, HSTS, X-Frame-Options, etc.)
- Cookie flags check (`Secure`, `HttpOnly`, `SameSite`)
- Resource sizes & heavy files list
- CSV/JSON exports
- Single-file HTML report (with embedded screenshot & results)
- Batch scan from TXT of domains
- Parallel/concurrent scans
- Colored terminal output

---

## 📦 Installation

Clone and build manually:

```bash
git clone https://github.com/FatihZor/WebProbeGo.git
cd WebProbeGo
go mod tidy
go build -o webprobego ./cmd/webprobego
```

---

## 🚀 Usage

```bash
# Basic keyword search
webprobego --domain=example.com --find=portfolio

# Search using a list of keywords from a file
webprobego --domain=example.com --find-file=terms.txt

# Extract meta info
webprobego --domain=example.com --meta

# Collect link inventory
webprobego --domain=example.com --links

# Capture network requests and take a screenshot
webprobego --domain=example.com --network --screenshot=site.png
```

