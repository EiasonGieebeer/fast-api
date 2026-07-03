package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

const (
	upstreamLatestReleaseURL = "https://api.github.com/repos/QuantumNous/new-api/releases/latest"
	maxReleaseResponseBytes  = 1 << 20
)

type upstreamRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

func GetLatestUpstreamRelease(c *gin.Context) {
	client := buildRatioSyncHTTPClient(strings.TrimSpace(common.GetEnvOrDefaultString("SYNC_HTTP_PROXY", "")))
	if client == nil {
		client = &http.Client{}
	}
	client.Timeout = 15 * time.Second

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, upstreamLatestReleaseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to create GitHub request"})
		return
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "fast-api-update-checker")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": fmt.Sprintf("failed to contact GitHub Releases API: %v", err)})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": fmt.Sprintf("GitHub Releases API returned HTTP %d", resp.StatusCode)})
		return
	}

	var release upstreamRelease
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxReleaseResponseBytes))
	if err := decoder.Decode(&release); err != nil || release.TagName == "" {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "GitHub Releases API returned an invalid response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": release})
}
