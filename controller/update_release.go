package controller

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

const (
	upstreamLatestReleaseURL     = "https://api.github.com/repos/QuantumNous/new-api/releases/latest"
	upstreamLatestReleaseFeedURL = "https://github.com/QuantumNous/new-api/releases.atom"
	maxReleaseResponseBytes      = 1 << 20
	latestReleaseCacheTTL        = 5 * time.Minute
)

type upstreamRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

type upstreamReleaseFeed struct {
	Entries []upstreamReleaseFeedEntry `xml:"entry"`
}

type upstreamReleaseFeedEntry struct {
	Title   string                    `xml:"title"`
	Updated string                    `xml:"updated"`
	Links   []upstreamReleaseFeedLink `xml:"link"`
}

type upstreamReleaseFeedLink struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

var latestUpstreamReleaseCache = struct {
	sync.RWMutex
	release   upstreamRelease
	expiresAt time.Time
}{}

func GetLatestUpstreamRelease(c *gin.Context) {
	latestUpstreamReleaseCache.RLock()
	cachedRelease := latestUpstreamReleaseCache.release
	cacheIsFresh := cachedRelease.TagName != "" && time.Now().Before(latestUpstreamReleaseCache.expiresAt)
	latestUpstreamReleaseCache.RUnlock()
	if cacheIsFresh {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": cachedRelease})
		return
	}

	clients := make([]*http.Client, 0, 2)
	directClient := buildRatioSyncHTTPClient("")
	if directClient == nil {
		directClient = &http.Client{}
	}
	directClient.Timeout = 15 * time.Second
	clients = append(clients, directClient)

	proxyURL := strings.TrimSpace(common.GetEnvOrDefaultString("SYNC_HTTP_PROXY", ""))
	if proxyURL != "" {
		proxyClient := buildRatioSyncHTTPClient(proxyURL)
		if proxyClient != nil {
			proxyClient.Timeout = 15 * time.Second
			clients = append(clients, proxyClient)
		}
	}

	release, err := fetchLatestUpstreamRelease(
		c.Request.Context(),
		upstreamLatestReleaseURL,
		upstreamLatestReleaseFeedURL,
		strings.TrimSpace(common.GetEnvOrDefaultString("GITHUB_API_TOKEN", "")),
		clients...,
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": err.Error()})
		return
	}

	latestUpstreamReleaseCache.Lock()
	latestUpstreamReleaseCache.release = release
	latestUpstreamReleaseCache.expiresAt = time.Now().Add(latestReleaseCacheTTL)
	latestUpstreamReleaseCache.Unlock()

	c.JSON(http.StatusOK, gin.H{"success": true, "data": release})
}

func fetchLatestUpstreamRelease(
	ctx context.Context,
	apiURL string,
	feedURL string,
	token string,
	clients ...*http.Client,
) (upstreamRelease, error) {
	if len(clients) == 0 {
		clients = append(clients, &http.Client{Timeout: 15 * time.Second})
	}

	attemptErrors := make([]string, 0, len(clients)*2)
	for _, client := range clients {
		release, err := fetchLatestUpstreamReleaseFromAPI(ctx, client, apiURL, token)
		if err == nil {
			return release, nil
		}
		attemptErrors = append(attemptErrors, err.Error())
	}

	for _, client := range clients {
		release, err := fetchLatestUpstreamReleaseFromFeed(ctx, client, feedURL)
		if err == nil {
			return release, nil
		}
		attemptErrors = append(attemptErrors, err.Error())
	}

	return upstreamRelease{}, fmt.Errorf("failed to check upstream release: %s", strings.Join(attemptErrors, "; "))
}

func fetchLatestUpstreamReleaseFromAPI(
	ctx context.Context,
	client *http.Client,
	apiURL string,
	token string,
) (upstreamRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return upstreamRelease{}, fmt.Errorf("failed to create GitHub API request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "fast-api-update-checker")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases API returned HTTP %d", resp.StatusCode)
	}

	var release upstreamRelease
	if err := common.DecodeJson(io.LimitReader(resp.Body, maxReleaseResponseBytes), &release); err != nil {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases API returned invalid JSON: %w", err)
	}
	if strings.TrimSpace(release.TagName) == "" {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases API returned a release without a tag")
	}
	return release, nil
}

func fetchLatestUpstreamReleaseFromFeed(
	ctx context.Context,
	client *http.Client,
	feedURL string,
) (upstreamRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return upstreamRelease{}, fmt.Errorf("failed to create GitHub Releases feed request: %w", err)
	}
	req.Header.Set("Accept", "application/atom+xml")
	req.Header.Set("User-Agent", "fast-api-update-checker")

	resp, err := client.Do(req)
	if err != nil {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases feed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases feed returned HTTP %d", resp.StatusCode)
	}

	var feed upstreamReleaseFeed
	decoder := xml.NewDecoder(io.LimitReader(resp.Body, maxReleaseResponseBytes))
	if err := decoder.Decode(&feed); err != nil {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases feed returned invalid XML: %w", err)
	}
	if len(feed.Entries) == 0 || strings.TrimSpace(feed.Entries[0].Title) == "" {
		return upstreamRelease{}, fmt.Errorf("GitHub Releases feed did not contain a release")
	}

	entry := feed.Entries[0]
	release := upstreamRelease{
		TagName:     strings.TrimSpace(entry.Title),
		Name:        strings.TrimSpace(entry.Title),
		PublishedAt: strings.TrimSpace(entry.Updated),
	}
	for _, link := range entry.Links {
		if link.Rel == "" || link.Rel == "alternate" {
			release.HTMLURL = strings.TrimSpace(link.Href)
			if release.HTMLURL != "" {
				break
			}
		}
	}
	return release, nil
}
