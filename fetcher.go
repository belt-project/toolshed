package toolshed

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type fetcher interface {
	Fetch(string) (string, error)
	Invalidate()
}

type githubFetcher struct {
	logger *log.Logger
	repo   string

	mu    sync.Mutex
	cache map[string]string
}

func (g *githubFetcher) Fetch(version string) (string, error) {
	script := g.cacheGet(version)

	if script != "" {
		return script, nil
	}

	g.logger.Printf("not cached, fetching script from github...")

	url := g.generateURL(version)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetcher request failed for: %s", url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	script = string(body)
	g.cachePut(version, script)

	return script, err
}

func (g *githubFetcher) Invalidate() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache = make(map[string]string)
}

func (g *githubFetcher) cacheGet(version string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.cache[version]
}

func (g *githubFetcher) cachePut(version, script string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache[version] = script
}

func (g *githubFetcher) generateURL(version string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/setup.sh", g.repo, version)
}
