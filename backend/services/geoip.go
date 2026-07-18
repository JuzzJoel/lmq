package services

import (
	"log"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

// GeoIPService wraps the maxminddb reader for country lookups.
// It gracefully handles a missing .mmdb file by returning "XX" for all lookups.
type GeoIPService struct {
	reader *maxminddb.Reader
}

// mmdbResult is a struct decoding city, subdivision (region), and country
// from the MaxMind City database.
type mmdbResult struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Subdivisions []struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// NewGeoIPService attempts to load the MaxMind GeoLite2-City DB from the given path.
// If the DB file is missing or unreadable, it logs a warning and returns a service
// with a nil reader that falls back to default unknown values for all lookups.
func NewGeoIPService(dbPath string) *GeoIPService {
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Printf("[GeoIPService]: Warning — Failed to open GeoIP DB at %s: %v. Lookups will return 'Unknown'/'XX'.", dbPath, err)
		return &GeoIPService{reader: nil}
	}
	log.Printf("[GeoIPService]: Loaded GeoLite2-City DB from %s", dbPath)
	return &GeoIPService{reader: db}
}

// LookupLocation returns the city, region, and ISO 3166-1 alpha-2 country code for the given IP address.
// Returns "Unknown", "Unknown", "XX" if the reader is nil, the IP is invalid, or the lookup fails.
func (s *GeoIPService) LookupLocation(ipStr string) (city, region, countryCode string) {
	city = "Unknown"
	region = "Unknown"
	countryCode = "XX"

	if s.reader == nil {
		return
	}

	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return
	}

	var result mmdbResult
	if err := s.reader.Lookup(addr).Decode(&result); err != nil {
		return
	}

	if result.City.Names != nil && result.City.Names["en"] != "" {
		city = result.City.Names["en"]
	}
	if len(result.Subdivisions) > 0 && result.Subdivisions[0].Names != nil && result.Subdivisions[0].Names["en"] != "" {
		region = result.Subdivisions[0].Names["en"]
	}
	if result.Country.ISOCode != "" {
		countryCode = result.Country.ISOCode
	}

	return
}

// Close closes the underlying maxminddb reader if it exists.
func (s *GeoIPService) Close() {
	if s.reader != nil {
		s.reader.Close()
	}
}
