// Package news provides functionality for fetching and managing the launcher news feed.
package news

import (
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"hytale-launcher/internal/endpoints"
	"hytale-launcher/internal/ioutil"
)

// cacheDuration is the time between feed refreshes.
const cacheDuration = 30 * time.Minute

// Article represents a single news article in the feed.
type Article struct {
	// ID is the unique identifier for the article.
	ID string `json:"id"`

	// Title is the article headline.
	Title string `json:"title"`

	// Summary is a brief description of the article.
	Summary string `json:"summary"`

	// ImageURL is the URL to the article's thumbnail image.
	ImageURL string `json:"image_url"`

	// LinkURL is the URL to the full article.
	LinkURL string `json:"link_url"`

	// PublishedAt is the publication timestamp.
	PublishedAt string `json:"published_at"`
}

// feedResponse is the JSON structure returned by the feed endpoint.
type feedResponse struct {
	Articles []Article `json:"articles"`
}

var (
	// mu protects access to the cached feed data.
	mu sync.RWMutex

	// cachedArticles holds the most recently fetched articles.
	cachedArticles []Article

	// lastFetch is the timestamp of the last successful fetch.
	lastFetch time.Time

	// baseURL is the parsed base URL for resolving relative URLs.
	baseURL *url.URL
)

// GetFeedArticles fetches news articles and returns whether new ones are available.
// If forceRefresh is true, the cache is bypassed and fresh data is fetched.
// Otherwise, cached data is used if it's still valid.
func GetFeedArticles(forceRefresh bool) (bool, error) {
	mu.RLock()
	timeSinceLastFetch := time.Since(lastFetch)
	previousCount := len(cachedArticles)
	mu.RUnlock()

	// Return cached state if still fresh and not forcing refresh
	if timeSinceLastFetch < cacheDuration && !forceRefresh {
		return false, nil
	}

	// Acquire write lock for refresh
	mu.Lock()
	defer mu.Unlock()

	// Double-check after acquiring write lock
	if time.Since(lastFetch) < cacheDuration && !forceRefresh {
		return false, nil
	}

	// Fetch fresh data
	articles, err := fetch()
	if err != nil {
		slog.Error("failed to fetch news feed", "error", err)
		return false, err
	}

	// Check if there are new articles
	hasNew := len(articles) > previousCount

	// Update cache
	cachedArticles = articles
	lastFetch = time.Now()

	return hasNew, nil
}

// GetCachedArticles returns the current cached list of articles.
func GetCachedArticles() []Article {
	mu.RLock()
	defer mu.RUnlock()
	return cachedArticles
}

// fetch retrieves the news feed from the server.
func fetch() ([]Article, error) {
	feedURL := endpoints.Feed()

	response, err := ioutil.Get[feedResponse](http.DefaultClient, feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Parse and cache the base URL for resolving relative URLs
	base, err := url.Parse(endpoints.FeedBase())
	if err != nil {
		slog.Error("failed to parse feed base URL", "error", err)
	} else {
		baseURL = base
	}

	// Resolve relative URLs in articles
	for i := range response.Articles {
		if response.Articles[i].ImageURL != "" {
			response.Articles[i].ImageURL = resolveURL(response.Articles[i].ImageURL)
		}
		if response.Articles[i].LinkURL != "" {
			response.Articles[i].LinkURL = resolveURL(response.Articles[i].LinkURL)
		}
	}

	return response.Articles, nil
}

// resolveURL resolves a potentially relative URL against the feed base URL.
// If the URL is already absolute or parsing fails, the original URL is returned.
func resolveURL(rawURL string) string {
	if baseURL == nil {
		return rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// If the URL is already absolute, return as-is
	if parsed.IsAbs() {
		return rawURL
	}

	// Resolve relative URL against base
	resolved := baseURL.ResolveReference(parsed)
	return resolved.String()
}

// ClearCache clears the cached feed data, forcing a refresh on the next call.
func ClearCache() {
	mu.Lock()
	defer mu.Unlock()
	cachedArticles = nil
	lastFetch = time.Time{}
}
