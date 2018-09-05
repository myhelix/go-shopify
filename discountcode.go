package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const discountCodesBasePath = "admin/price_rules"

// DiscountCodeService is an interface for interfacing with the discount code endpoints
// of the Shopify API.
// See: https://help.shopify.com/en/api/reference/discounts/discountcode
type DiscountCodeService interface {
	List(int64, interface{}) ([]DiscountCode, error)
	Get(int64, int64, interface{}) (*DiscountCode, error)
	Create(int64, DiscountCode) (*DiscountCode, error)
	Update(int64, int64, DiscountCode) (*DiscountCode, error)
	Delete(int64, int64) error
}

// DiscountCodeServiceOp handles communication with the discount code related methods of
// the Shopify API.
type DiscountCodeServiceOp struct {
	client *Client
}

type DiscountCode struct {
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Code        string           `json:"code,omitempty"`
	Type        string           `json:"type,omitempty"`
	CreatedAt   *time.Time       `json:"created_at,omitempty"`
	UpdatedAt   *time.Time       `json:"updated_at,omitempty"`
	ID          int64            `json:"id,omitempty"`
	PriceRuleID int64            `json:"price_rule_id,omitempty"`
	UsageCount  int              `json:"usage_count,omitempty"`
}

// Represents the result from the discount_codes/X.json endpoint
type DiscountCodeResource struct {
	DiscountCode *DiscountCode `json:"discount_code"`
}

// Represents the result from the discount_codes.json endpoint
type DiscountCodesResource struct {
	DiscountCodes []DiscountCode `json:"discount_codes"`
}

// List discount codes
func (s *DiscountCodeServiceOp) List(priceRuleID int64, options interface{}) ([]DiscountCode, error) {
	path := fmt.Sprintf("%s/%d/discount_codes.json", discountCodesBasePath, priceRuleID)
	resource := new(DiscountCodesResource)
	err := s.client.Get(path, resource, options)
	return resource.DiscountCodes, err
}

// Get discount code
func (s *DiscountCodeServiceOp) Get(priceRuleID int64, discountCodeID int64, options interface{}) (*DiscountCode, error) {
	path := fmt.Sprintf("%s/%d/discount_codes/%d.json", discountCodesBasePath, priceRuleID, discountCodeID)
	resource := new(DiscountCodeResource)
	err := s.client.Get(path, resource, options)
	return resource.DiscountCode, err
}

// Create a new discount code
func (s *DiscountCodeServiceOp) Create(priceRuleID int64, discountCode DiscountCode) (*DiscountCode, error) {
	path := fmt.Sprintf("%s/%d/discount_codes.json", discountCodesBasePath, priceRuleID)
	wrappedData := DiscountCodeResource{DiscountCode: &discountCode}
	resource := new(DiscountCodeResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.DiscountCode, err
}

// Update an existing discount code
func (s *DiscountCodeServiceOp) Update(priceRuleID int64, discountCodeID int64, discountCode DiscountCode) (*DiscountCode, error) {
	path := fmt.Sprintf("%s/%d/discount_codes/%d.json", discountCodesBasePath, priceRuleID, discountCodeID)
	wrappedData := DiscountCodeResource{DiscountCode: &discountCode}
	resource := new(DiscountCodeResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.DiscountCode, err
}

// Delete an existing discount code
func (s *DiscountCodeServiceOp) Delete(priceRuleID int64, discountCodeID int64) error {
	path := fmt.Sprintf("%s/%d/discount_codes/%d.json", discountCodesBasePath, priceRuleID, discountCodeID)
	return s.client.Delete(path)
}
