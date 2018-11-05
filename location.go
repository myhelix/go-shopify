package goshopify

import (
	"fmt"
	"time"
)

const locationsBasePath = "admin/locations"

// LocationService is an interface for interfacing with the location endpoints
// of the Shopify API.
// https://help.shopify.com/en/api/reference/inventory/location
type LocationService interface {
	List(interface{}) ([]Location, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Location, error)
}

// LocationServiceOp handles communication with the location related methods of
// the Shopify API.
type LocationServiceOp struct {
	client *Client
}

// Location represents a Shopify location.
type Location struct {
	ID           int64      `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	Address1     string     `json:"address1,omitempty"`
	Address2     string     `json:"address2,omitempty"`
	City         string     `json:"city,omitempty"`
	Zip          string     `json:"zip,omitempty"`
	Province     string     `json:"province,omitempty"`
	Country      string     `json:"country,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	CountryCode  string     `json:"country_code,omitempty"`
	CountryName  string     `json:"country_name,omitempty"`
	ProvinceCode string     `json:"province_code,omitempty"`
	Legacy       bool       `json:"legacy,omitempty"`
	Active       bool       `json:"active,omitempty"`
}

// Represents the result from the location/X.json endpoint.
type LocationResource struct {
	Location *Location `json:"location"`
}

// Represents the result from the locations/X.json endpoint.
type LocationsResource struct {
	Locations []Location `json:"locations"`
}

// List locations.
func (s *LocationServiceOp) List(options interface{}) ([]Location, error) {
	path := fmt.Sprintf("%s.json", locationsBasePath)
	resource := new(LocationsResource)
	err := s.client.Get(path, resource, options)
	return resource.Locations, err
}

// Count locations.
func (s *LocationServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", locationsBasePath)
	return s.client.Count(path, options)
}

// Get location.
func (s *LocationServiceOp) Get(locationID int64, options interface{}) (*Location, error) {
	path := fmt.Sprintf("%s/%v.json", locationsBasePath, locationID)
	resource := new(LocationResource)
	err := s.client.Get(path, resource, options)
	return resource.Location, err
}
