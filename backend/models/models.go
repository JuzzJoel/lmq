package models

import "time"

// Link represents a shortened URL entry in the database.
type Link struct {
	ID         int64      `json:"id"`
	Token      string     `json:"token"`
	LongURL    string     `json:"long_url"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	ClickCount   int64      `json:"click_count"`
	HasPassword      bool        `json:"has_password"`
	PasswordHash     *string     `json:"-"` // Not exported in JSON
	BurnAfterReading bool        `json:"burn_after_reading"`
	Routes           []RouteSpec `json:"routes,omitempty"`
}

// ClickEvent represents an analytics record for a single link visit.
type ClickEvent struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IPAddress   string    `json:"ip_address,omitempty"`
	City        string    `json:"city"`
	Region      string    `json:"region"`
	CountryCode string    `json:"country_code"`
	UserAgent string    `json:"user_agent,omitempty"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	IsMobile  bool      `json:"is_mobile"`
	Referer   string    `json:"referer,omitempty"`
}

// RouteSpec defines a single A/B testing destination with a weight percentage.
type RouteSpec struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

// ShortenRequest defines the expected JSON payload for creating a new short link.
type ShortenRequest struct {
	URL         string     `json:"url"`
	CustomToken string     `json:"custom_token,omitempty"`
	ExpiresIn   int        `json:"expires_in,omitempty"` // in hours
	Password         string     `json:"password,omitempty"`
	Routes           []RouteSpec `json:"routes,omitempty"`
	BurnAfterReading bool        `json:"burn_after_reading,omitempty"`
}

// APIResponse is a generic response wrapper for API responses.
type APIResponse[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

// LinkAnalytics aggregates statistics for a specific short link.
type LinkAnalytics struct {
	Token        string         `json:"token"`
	LongURL      string         `json:"long_url"`
	TotalClicks  int64          `json:"total_clicks"`
	ClicksByDay   []DayCount     `json:"clicks_by_day"`
	Cities        []CityCount    `json:"cities"`
	Regions       []RegionCount  `json:"regions"`
	CountryGroups []CountryCount `json:"country_groups"`
	Browsers     []BrowserCount `json:"browsers"`
	RecentClicks []ClickEvent   `json:"recent_clicks"`
	Routes           []RouteSpec `json:"routes,omitempty"`
	BurnAfterReading bool        `json:"burn_after_reading,omitempty"`
}

// DayCount represents click counts grouped by day.
type DayCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// CityCount represents click counts grouped by city.
type CityCount struct {
	City  string `json:"city"`
	Count int64  `json:"count"`
}

// RegionCount represents click counts grouped by region.
type RegionCount struct {
	Region string `json:"region"`
	Count  int64  `json:"count"`
}

// CountryCount represents click counts grouped by country code.
type CountryCount struct {
	CountryCode string `json:"country_code"`
	Count       int64  `json:"count"`
}

// BrowserCount represents click counts grouped by browser.
type BrowserCount struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}
