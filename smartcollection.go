package goshopify

import (
	"fmt"
	"net/http"
	"time"
)

const smartCollectionsBasePath = "smart_collections"
const smartCollectionsResourceName = "collections"

// SmartCollectionService is an interface for interacting with the smart
// collection endpoints of the Shopify API.
// See https://help.shopify.com/api/reference/smartcollection
type SmartCollectionService interface {
	List(interface{}) ([]SmartCollection, error)
	ListWithPagination(interface{}) ([]SmartCollection, *Pagination, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*SmartCollection, error)
	Create(SmartCollection) (*SmartCollection, error)
	Update(SmartCollection) (*SmartCollection, error)
	Delete(int64) error

	// MetafieldsService used for SmartCollection resource to communicate with Metafields resource
	MetafieldsService
}

// SmartCollectionServiceOp handles communication with the smart collection
// related methods of the Shopify API.
type SmartCollectionServiceOp struct {
	client *Client
}

type Rule struct {
	Column    string `json:"column,omitempty"`
	Relation  string `json:"relation,omitempty"`
	Condition string `json:"condition,omitempty"`
}

// SmartCollection represents a Shopify smart collection.
type SmartCollection struct {
	ID             int64       `json:"id,omitempty"`
	Handle         string      `json:"handle,omitempty"`
	Title          string      `json:"title,omitempty"`
	UpdatedAt      *time.Time  `json:"updated_at,omitempty"`
	BodyHTML       string      `json:"body_html,omitempty"`
	SortOrder      string      `json:"sort_order,omitempty"`
	TemplateSuffix string      `json:"template_suffix,omitempty"`
	Image          Image       `json:"image,omitempty"`
	Published      bool        `json:"published,omitempty"`
	PublishedAt    *time.Time  `json:"published_at,omitempty"`
	PublishedScope string      `json:"published_scope,omitempty"`
	Rules          []Rule      `json:"rules,omitempty"`
	Disjunctive    bool        `json:"disjunctive,omitempty"`
	Metafields     []Metafield `json:"metafields,omitempty"`
}

// SmartCollectionResource represents the result from the smart_collections/X.json endpoint
type SmartCollectionResource struct {
	Collection *SmartCollection `json:"smart_collection"`
}

// SmartCollectionsResource represents the result from the smart_collections.json endpoint
type SmartCollectionsResource struct {
	Collections []SmartCollection `json:"smart_collections"`
}

// List smart collections
func (s *SmartCollectionServiceOp) List(options interface{}) ([]SmartCollection, error) {
	smartCollections, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return smartCollections, nil
}

// List smart collections with pagination
func (s *SmartCollectionServiceOp) ListWithPagination(options interface{}) ([]SmartCollection, *Pagination, error) {
	path := fmt.Sprintf("%s.json", smartCollectionsBasePath)
	resource := new(SmartCollectionsResource)
	headers := http.Header{}

	headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
	if err != nil {
		return nil, nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource.Collections, pagination, nil
}

// Count smart collections
func (s *SmartCollectionServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", smartCollectionsBasePath)
	return s.client.Count(path, options)
}

// Get individual smart collection
func (s *SmartCollectionServiceOp) Get(collectionID int64, options interface{}) (*SmartCollection, error) {
	path := fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collectionID)
	resource := new(SmartCollectionResource)
	err := s.client.Get(path, resource, options)
	return resource.Collection, err
}

// Create a new smart collection
// See Image for the details of the Image creation for a collection.
func (s *SmartCollectionServiceOp) Create(collection SmartCollection) (*SmartCollection, error) {
	path := fmt.Sprintf("%s.json", smartCollectionsBasePath)
	wrappedData := SmartCollectionResource{Collection: &collection}
	resource := new(SmartCollectionResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Collection, err
}

// Update an existing smart collection
func (s *SmartCollectionServiceOp) Update(collection SmartCollection) (*SmartCollection, error) {
	path := fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collection.ID)
	wrappedData := SmartCollectionResource{Collection: &collection}
	resource := new(SmartCollectionResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Collection, err
}

// Delete an existing smart collection.
func (s *SmartCollectionServiceOp) Delete(collectionID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", smartCollectionsBasePath, collectionID))
}

// List metafields for a smart collection
func (s *SmartCollectionServiceOp) ListMetafields(smartCollectionID int64, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.List(options)
}

// Count metafields for a smart collection
func (s *SmartCollectionServiceOp) CountMetafields(smartCollectionID int64, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.Count(options)
}

// Get individual metafield for a smart collection
func (s *SmartCollectionServiceOp) GetMetafield(smartCollectionID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.Get(metafieldID, options)
}

// Create a new metafield for a smart collection
func (s *SmartCollectionServiceOp) CreateMetafield(smartCollectionID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.Create(metafield)
}

// Update an existing metafield for a smart collection
func (s *SmartCollectionServiceOp) UpdateMetafield(smartCollectionID int64, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.Update(metafield)
}

// // Delete an existing metafield for a smart collection
func (s *SmartCollectionServiceOp) DeleteMetafield(smartCollectionID int64, metafieldID int64) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: smartCollectionsResourceName, resourceID: smartCollectionID}
	return metafieldService.Delete(metafieldID)
}
