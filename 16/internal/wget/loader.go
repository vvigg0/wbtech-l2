package wget

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type loader struct {
	rootDir string
	client  *http.Client
}

func newLoader(cfg *Config) (*loader, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	rootDir := filepath.Join(cfg.RootPath, u.Host)
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, err
	}

	return &loader{
		rootDir: rootDir,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (l *loader) fetch(rawURL string) ([]byte, string, error) {
	resp, err := l.client.Get(rawURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return body, resp.Header.Get("Content-Type"), nil
}

func (l *loader) save(rel string, body []byte) error {
	fullPath := filepath.Join(l.rootDir, filepath.FromSlash(rel))

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(fullPath, body, 0644); err != nil {
		return err
	}

	fmt.Printf("Сохранено: %s\n", fullPath)
	return nil
}

func isSkippableLink(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return true
	}

	low := strings.ToLower(raw)
	return strings.HasPrefix(low, "#") ||
		strings.HasPrefix(low, "mailto:") ||
		strings.HasPrefix(low, "tel:") ||
		strings.HasPrefix(low, "javascript:") ||
		strings.HasPrefix(low, "data:")
}

func fullURL(base *url.URL, raw string, domain string) (*url.URL, bool) {
	raw = strings.TrimSpace(raw)
	if isSkippableLink(raw) {
		return nil, false
	}

	if strings.HasPrefix(raw, "//") {
		raw = base.Scheme + ":" + raw
	}

	ref, err := url.Parse(raw)
	if err != nil {
		return nil, false
	}

	abs := base.ResolveReference(ref)
	abs.Fragment = ""
	abs.RawQuery = ""
	abs.ForceQuery = false

	if abs.Host != "" && abs.Host != domain {
		return nil, false
	}

	if abs.Path == "" {
		abs.Path = "/"
	}

	return abs, true
}

func localRelPath(u *url.URL) string {
	p := u.Path
	if p == "" || p == "/" {
		return "index.html"
	}

	p = path.Clean(p)

	if strings.HasSuffix(u.Path, "/") {
		p = path.Join(p, "index.html")
	} else if path.Ext(p) == "" {
		p = path.Join(p, "index.html")
	}

	return strings.TrimPrefix(p, "/")
}

func makeRel(fromFileRel, toFileRel string) string {
	fromDir := filepath.Dir(filepath.FromSlash(fromFileRel))
	to := filepath.FromSlash(toFileRel)

	rel, err := filepath.Rel(fromDir, to)
	if err != nil {
		return toFileRel
	}

	return filepath.ToSlash(rel)
}

func isHTMLContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "text/html")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
