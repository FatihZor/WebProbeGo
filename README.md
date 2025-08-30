# WebProbeGo

**WebProbeGo** — a fast and versatile web analysis tool built with Go and Chromedp.  
It allows you to probe websites, capture screenshots, inspect network requests, and search for keywords — all from a simple CLI.

---

## ✨ Features

- **Network request logging** (`--network`)  
  Capture all requests a page makes and report their HTTP status codes.

- **Keyword/phrase search** (`--find`, `--find-file`)  
  Search for a single term or a list of terms (from a `.txt` file) in the page body.

- **Screenshots** (`--screenshot`)  
  Save a full-page PNG screenshot of the target site.

- **Execution time measurement**  
  Report how long the operation took in seconds.

---

## 🚧 Roadmap

Planned improvements for upcoming versions:

- Extract meta information: `<title>`, `<meta description>`, `<meta keywords>`, favicon.
- Heading structure map (H1–H3, optionally deeper).
- Full link inventory: internal vs external references.
- Broken link checker for `<a href>` and `<img src>`.
- SSL/TLS certificate inspection (validity & expiration date).
- Redirect chain tracking and loop detection.
- Selector-only screenshots (e.g. `#main-content`).
- Mobile & tablet viewport emulation (device profiles).
- Security headers report: CSP, HSTS, X-Frame-Options, etc.
- Cookie security flags check (`Secure`, `HttpOnly`, `SameSite`).
- Resource size reporting (large JS/CSS/images flagged).
- CSV/JSON export for structured analysis results.
- Single-file HTML report (with embedded screenshot & results).
- Batch scanning multiple domains from a TXT file.
- Parallel/concurrent scanning for faster batch runs.
- Colored terminal output for better readability.

---

## 📦 Installation

Clone and build manually:

```bash
git clone https://github.com/FatihZor/WebProbeGo.git
cd WebProbeGo
go build -o webprobego ./cmd/webprobego
```