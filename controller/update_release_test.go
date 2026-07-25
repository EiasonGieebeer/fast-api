package controller

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type releaseRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn releaseRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func TestFetchLatestUpstreamReleaseRetriesWithNextClientAfterForbidden(t *testing.T) {
	firstCalls := 0
	secondCalls := 0
	firstClient := &http.Client{Transport: releaseRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		firstCalls++
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"message":"API rate limit exceeded"}`)),
			Request:    request,
		}, nil
	})}
	secondClient := &http.Client{Transport: releaseRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		secondCalls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"tag_name":"v1.0.0-rc.21",
				"name":"v1.0.0-rc.21",
				"html_url":"https://github.com/QuantumNous/new-api/releases/tag/v1.0.0-rc.21",
				"published_at":"2026-07-11T15:01:26Z"
			}`)),
			Request: request,
		}, nil
	})}

	release, err := fetchLatestUpstreamRelease(
		context.Background(),
		"https://api.example.test/releases/latest",
		"https://example.test/releases.atom",
		"",
		firstClient,
		secondClient,
	)

	require.NoError(t, err)
	assert.Equal(t, "v1.0.0-rc.21", release.TagName)
	assert.Equal(t, 1, firstCalls)
	assert.Equal(t, 1, secondCalls)
}

func TestFetchLatestUpstreamReleaseFallsBackToAtomFeed(t *testing.T) {
	client := &http.Client{Transport: releaseRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if strings.HasSuffix(request.URL.Path, ".atom") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>
					<feed xmlns="http://www.w3.org/2005/Atom">
						<entry>
							<title>v1.0.0-rc.21</title>
							<updated>2026-07-11T15:01:26Z</updated>
							<link rel="alternate" href="https://github.com/QuantumNous/new-api/releases/tag/v1.0.0-rc.21"/>
						</entry>
					</feed>`)),
				Request: request,
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"message":"API rate limit exceeded"}`)),
			Request:    request,
		}, nil
	})}

	release, err := fetchLatestUpstreamRelease(
		context.Background(),
		"https://api.example.test/releases/latest",
		"https://example.test/releases.atom",
		"",
		client,
	)

	require.NoError(t, err)
	assert.Equal(t, "v1.0.0-rc.21", release.TagName)
	assert.Equal(t, "v1.0.0-rc.21", release.Name)
	assert.Equal(t, "2026-07-11T15:01:26Z", release.PublishedAt)
	assert.Equal(t, "https://github.com/QuantumNous/new-api/releases/tag/v1.0.0-rc.21", release.HTMLURL)
}

func TestFetchLatestUpstreamReleaseSendsConfiguredGitHubToken(t *testing.T) {
	client := &http.Client{Transport: releaseRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, "Bearer test-token", request.Header.Get("Authorization"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v1.0.0-rc.21"}`)),
			Request:    request,
		}, nil
	})}

	_, err := fetchLatestUpstreamRelease(
		context.Background(),
		"https://api.example.test/releases/latest",
		"https://example.test/releases.atom",
		"test-token",
		client,
	)

	require.NoError(t, err)
}
