package goshopify

import (
	"fmt"
	"time"
)

const productsBasePath = "admin/products"
const productsResourceName = "products"

// ProductService is an interface for interfacing with the product endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/product
type ProductService interface {
	List(interface{}) ([]Product, error)
	Count(interface{}) (int, error)
	Get(int, interface{}) (*Product, error)
	Create(Product) (*Product, error)
	Update(Product) (*Product, error)
	Delete(int) error

	// MetafieldsService used for Product resource to communicate with Metafields resource
	MetafieldsService
}

// ProductServiceOp handles communication with the product related methods of
// the Shopify API.
type ProductServiceOp struct {
	client *Client
}

// Product represents a Shopify product
type Product struct {
	ID                             int             `json:"id"`
	Title                          string          `json:"title"`
	BodyHTML                       string          `json:"body_html"`
	Vendor                         string          `json:"vendor"`
	ProductType                    string          `json:"product_type"`
	Handle                         string          `json:"handle"`
	CreatedAt                      *time.Time      `json:"created_at"`
	UpdatedAt                      *time.Time      `json:"updated_at"`
	PublishedAt                    *time.Time      `json:"published_at"`
	PublishedScope                 string          `json:"published_scope"`
	Tags                           string          `json:"tags"`
	Options                        []ProductOption `json:"options"`
	Variants                       []Variant       `json:"variants"`
	Image                          Image           `json:"image"`
	Images                         []Image         `json:"images"`
	TemplateSuffix                 string          `json:"template_suffix"`
	MetafieldsGlobalTitleTag       string          `json:"metafields_global_title_tag"`
	MetafieldsGlobalDescriptionTag string          `json:"metafields_global_description_tag"`
	Metafields                     []Metafield     `json:"metafields"`
}

// The options provided by Shopify
type ProductOption struct {
	ID        int      `json:"id"`
	ProductID int      `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

// Represents the result from the products/X.json endpoint
type ProductResource struct {
	Product *Product `json:"product"`
}

// Represents the result from the products.json endpoint
type ProductsResource struct {
	Products []Product `json:"products"`
}

// List products
func (s *ProductServiceOp) List(options interface{}) ([]Product, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	resource := new(ProductsResource)
	err := s.client.Get(path, resource, options)
	return resource.Products, err
}

// Count products
func (s *ProductServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", productsBasePath)
	return s.client.Count(path, options)
}

// Get individual product
func (s *ProductServiceOp) Get(productID int, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, productID)
	resource := new(ProductResource)
	err := s.client.Get(path, resource, options)
	return resource.Product, err
}

// Create a new product
func (s *ProductServiceOp) Create(product Product) (*Product, error) {
	path := fmt.Sprintf("%s.json", productsBasePath)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Product, err
}

// Update an existing product
func (s *ProductServiceOp) Update(product Product) (*Product, error) {
	path := fmt.Sprintf("%s/%d.json", productsBasePath, product.ID)
	wrappedData := ProductResource{Product: &product}
	resource := new(ProductResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Product, err
}

// Delete an existing product
func (s *ProductServiceOp) Delete(productID int) error {
	return s.client.Delete(fmt.Sprintf("%s/%d.json", productsBasePath, productID))
}

// List metafields for a product
func (s *ProductServiceOp) ListMetafields(productID int, options interface{}) ([]Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.List(options)
}

// Count metafields for a product
func (s *ProductServiceOp) CountMetafields(productID int, options interface{}) (int, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Count(options)
}

// Get individual metafield for a product
func (s *ProductServiceOp) GetMetafield(productID int, metafieldID int, options interface{}) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Get(metafieldID, options)
}

// Create a new metafield for a product
func (s *ProductServiceOp) CreateMetafield(productID int, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Create(metafield)
}

// Update an existing metafield for a product
func (s *ProductServiceOp) UpdateMetafield(productID int, metafield Metafield) (*Metafield, error) {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Update(metafield)
}

// // Delete an existing metafield for a product
func (s *ProductServiceOp) DeleteMetafield(productID int, metafieldID int) error {
	metafieldService := &MetafieldServiceOp{client: s.client, resource: productsResourceName, resourceID: productID}
	return metafieldService.Delete(metafieldID)
}
