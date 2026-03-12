package wget

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type crawler struct {
	domain string

	mu      sync.Mutex
	visited map[string]struct{}
	wg      sync.WaitGroup

	tasks chan task
	found chan task

	maxDepth int

	loader *loader
}

type task struct {
	url   string
	depth int
}

func newCrawler(cfg *Config) (*crawler, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	l, err := newLoader(cfg)
	if err != nil {
		return nil, err
	}

	return &crawler{
		domain:   u.Host,
		visited:  make(map[string]struct{}),
		tasks:    make(chan task),
		found:    make(chan task, max(128, cfg.Workers*32)),
		maxDepth: cfg.Depth,
		loader:   l,
	}, nil
}

func (c *crawler) enqueue(rawURL string, depth int) {
	if c.maxDepth != -1 && depth > c.maxDepth {
		return
	}

	norm, ok := normalizeTaskURL(rawURL, c.domain)
	if !ok {
		return
	}

	c.mu.Lock()
	if _, exists := c.visited[norm]; exists {
		c.mu.Unlock()
		return
	}
	c.visited[norm] = struct{}{}
	c.wg.Add(1)
	c.mu.Unlock()

	c.found <- task{url: norm, depth: depth}
}

func (c *crawler) worker() {
	for t := range c.tasks {
		func() {
			defer c.wg.Done()

			body, contentType, err := c.loader.fetch(t.url)
			if err != nil {
				fmt.Fprintln(os.Stderr, "fetch error:", err)
				return
			}

			u, err := url.Parse(t.url)
			if err != nil {
				fmt.Fprintln(os.Stderr, "url parse error:", err)
				return
			}

			currentRel := localRelPath(u)

			if isHTMLContentType(contentType) {
				rewritten, pageLinks, assetLinks, err := c.rewriteHTML(u, currentRel, body)
				if err != nil {
					fmt.Fprintln(os.Stderr, "html parse error:", err)
					return
				}

				if err := c.loader.save(currentRel, rewritten); err != nil {
					fmt.Fprintln(os.Stderr, "save html error:", err)
					return
				}

				for _, link := range assetLinks {
					c.enqueue(link, t.depth)
				}

				for _, link := range pageLinks {
					c.enqueue(link, t.depth+1)
				}

				return
			}

			if err := c.loader.save(currentRel, body); err != nil {
				fmt.Fprintln(os.Stderr, "save asset error:", err)
				return
			}
		}()
	}
}

func (c *crawler) dispatch() {
	queue := make([]task, 0)

	for c.found != nil || len(queue) > 0 {
		var (
			out  chan task
			next task
		)

		if len(queue) > 0 {
			out = c.tasks
			next = queue[0]
		}

		select {
		case t, ok := <-c.found:
			if !ok {
				c.found = nil
				continue
			}
			queue = append(queue, t)

		case out <- next:
			queue = queue[1:]
		}
	}

	close(c.tasks)
}

func (c *crawler) rewriteHTML(baseURL *url.URL, currentRel string, body []byte) ([]byte, []string, []string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, nil, nil, err
	}

	pageLinks := make([]string, 0)
	assetLinks := make([]string, 0)

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				c.rewriteAttr(n, "href", baseURL, currentRel, &pageLinks)
			case "link":
				c.rewriteAttr(n, "href", baseURL, currentRel, &assetLinks)
			case "script":
				c.rewriteAttr(n, "src", baseURL, currentRel, &assetLinks)
			case "img":
				c.rewriteAttr(n, "src", baseURL, currentRel, &assetLinks)
			case "source":
				c.rewriteAttr(n, "src", baseURL, currentRel, &assetLinks)
			case "iframe":
				c.rewriteAttr(n, "src", baseURL, currentRel, &assetLinks)
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, nil, nil, err
	}

	return buf.Bytes(), pageLinks, assetLinks, nil
}

func (c *crawler) rewriteAttr(n *html.Node, attrName string, baseURL *url.URL, currentRel string, collected *[]string) {
	for i := range n.Attr {
		if n.Attr[i].Key != attrName {
			continue
		}

		abs, ok := fullURL(baseURL, n.Attr[i].Val, c.domain)
		if !ok {
			return
		}

		targetRel := localRelPath(abs)
		n.Attr[i].Val = makeRel(currentRel, targetRel)
		*collected = append(*collected, abs.String())
		return
	}
}

func crawl(cfg *Config) error {
	c, err := newCrawler(cfg)
	if err != nil {
		return err
	}

	go c.dispatch()

	c.enqueue(cfg.URL, 0)

	var workersWG sync.WaitGroup
	workersWG.Add(cfg.Workers)

	for i := 0; i < cfg.Workers; i++ {
		go func() {
			defer workersWG.Done()
			c.worker()
		}()
	}

	go func() {
		c.wg.Wait()
		close(c.found)
	}()

	workersWG.Wait()
	return nil
}

func normalizeTaskURL(raw string, domain string) (string, bool) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", false
	}

	if u.Host != "" && u.Host != domain {
		return "", false
	}

	u.Fragment = ""
	u.RawQuery = ""
	u.ForceQuery = false

	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), true
}
