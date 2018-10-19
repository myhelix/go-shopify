package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const variantsBasePath = "admin/variants"
const variantsResourceName = "variants"

// VariantService is an interface for interacting with the variant endpoints
// of the Shopify API.
// See https://help.shopify.com/api/reference/product_variant
type VariantService interface {
	List(int64, interface{}) ([]Variant, error)
	Count(int64, interface{}) (int, error)
	Get(int64, interface{}) (*Variant, error)
	Create(int64, Variant) (*Variant, error)
	Update(Variant) (*Variant, error)
	Delete(int64, int64) error

	// MetafieldsService used for Variant resource to communicate with Metafields resource
	VariantMetafieldsService
}

// VariantServiceOp handles communication with the variant related methods of
// the Shopify API.
type VariantServiceOp struct {
	client *Client
}

// Variant represents a Shopify variant
type Variant struct {
	ID                   int64            `json:"id,omitempty"`
	ProductID            int64            `json:"product_id,omitempty"`
	Title                string           `json:"title,omitempty"`
	Sku                  string           `json:"sku,omitempty"`
	Position             int              `json:"position,omitempty"`
	Grams                int              `json:"grams,omitempty"`
	InventoryPolicy      string           `json:"inventory_policy,omitempty"`
	Price                *decimal.Decimal `json:"price,omitempty"`
	CompareAtPrice       *decimal.Decimal `json:"compare_at_price,omitempty"`
	FulfillmentService   string           `json:"fulfillment_service,omitempty"`
	InventoryManagement  string           `json:"inventory_management,omitempty"`
	Option1              string           `json:"option1,omitempty"`
	Option2              string           `json:"option2,omitempty"`
	Option3              string           `json:"option3,omitempty"`
	CreatedAt            *time.Time       `json:"created_at,omitempty"`
	UpdatedAt            *time.Time       `json:"updated_at,omitempty"`
	Taxable              bool             `json:"taxable,omitempty"`
	Barcode              string           `json:"barcode,omitempty"`
	ImageID              int64            `json:"image_id,omitempty"`
	InventoryQuantity    int              `json:"inventory_quantity,omitempty"`
	Weight               *decimal.Decimal `json:"weight,omitempty"`
	WeightUnit           string           `json:"weight_unit,omitempty"`
	OldInventoryQuantity int              `json:"old_inventory_quantity,omitempty"`
	RequireShipping      bool             `json:"requires_shipping,omitempty"`
	Metafields           []Metafield      `json:"metafields,omitempty"`
}

// VariantResource represents the result from the variants/X.json endpoint
type VariantResource struct {
	Variant *Variant `json:"variant"`
}

// VariantsResource represents the result from the products/X/variants.json endpoint
type VariantsResource struct {
	Variants []Variant `json:"variants"`
}

// List variants
func (s *VariantServiceOp) List(productID int64, options interface{}) ([]Variant, error) {
	path := fmt.Sprintf("%s/%d/variants.json", productsBasePath, productID)
	resource := new(VariantsResource)
	err := s.client.Get(path, resource, options)
	return resource.Variants, err
}

// Count variants
func (s *VariantServiceOp) Count(productID int64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%d/variants/count.json", productsBasePath, productID)
	return s.client.Count(path, options)
}

// Get individual variant
func (s *VariantServiceOp) Get(variantID int64, options interface{}) (*Variant, error) {
	path := fmt.Sprintf("%s/%d.json", variantsBasePath, variantID)
	resource := new(VariantResource)
	err := s.client.Get(path, resource, options)
	return resource.Variant, err
}

// Create a new variant
func (s *VariantServiceOp) Create(productID int64, variant Variant) (*Variant, error) {
	path := fmt.Sprintf("%s/%d/variants.json", productsBasePath, productID)
	wrappedData := VariantResource{Variant: &variant}
	resource := new(VariantResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Variant, err
}

// Update existing variant
func (s *VariantServiceOp) Update(variant Variant) (*Variant, error) {
	path := fmt.Sprintf("%s/%d.json", variantsBasePath, variant.ID)
	wrappedData := VariantResource{Variant: &variant}
	resource := new(VariantResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Variant, err
}

// Delete an existing variant
func (s *VariantServiceOp) Delete(productID int64, variantID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%d/variants/%d.json", productsBasePath, productID, variantID))
}

// List metafields for a variant
func (s *VariantServiceOp) ListMetafields(productID int64, variantID int64, options interface{}) ([]Metafield, error) {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.List(options)
}

// Count metafields for a variant
func (s *VariantServiceOp) CountMetafields(productID int64, variantID int64, options interface{}) (int, error) {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.Count(options)
}

// Get individual metafield for a variant
func (s *VariantServiceOp) GetMetafield(productID int64, variantID int64, metafieldID int64, options interface{}) (*Metafield, error) {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.Get(metafieldID, options)
}

// Create a new metafield for a variant
func (s *VariantServiceOp) CreateMetafield(productID int64, variantID int64, metafield Metafield) (*Metafield, error) {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.Create(metafield)
}

// Update an existing metafield for a variant
func (s *VariantServiceOp) UpdateMetafield(productID int64, variantID int64, metafield Metafield) (*Metafield, error) {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.Update(metafield)
}

// Delete an existing metafield for a variant
func (s *VariantServiceOp) DeleteMetafield(productID int64, variantID int64, metafieldID int64) error {
	variantMetafieldResource := fmt.Sprintf("products/%d/variants", productID)
	metafieldService := &MetafieldServiceOp{client: s.client, resource: variantMetafieldResource, resourceID: variantID}
	return metafieldService.Delete(metafieldID)
}
