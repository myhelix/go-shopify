package goshopify

import (
	"fmt"
	"time"
)

const priceRuleBasePath = "admin/price_rules"

// PriceRuleService is an interface for interfacing with the price rule endpoints
// of the Shopify API.
// https://help.shopify.com/en/api/reference/discounts/pricerule
type PriceRuleService interface {
	List(interface{}) ([]PriceRule, error)
	Get(int64, interface{}) (*PriceRule, error)
	Create(PriceRule) (*PriceRule, error)
	Update(PriceRule) (*PriceRule, error)
	Delete(int64) error
}

// PriceRuleServiceOp handles communication with the price rule related methods of
// the Shopify API.
type PriceRuleServiceOp struct {
	client *Client
}

type PriceRule struct {
	ID                                     int64                                   `json:"id,omitempty"`
	ValueType                              string                                  `json:"value_type,omitempty"`
	Value                                  string                                  `json:"value,omitempty"`
	CustomerSelection                      string                                  `json:"customer_selection,omitempty"`
	TargetType                             string                                  `json:"target_type,omitempty"`
	TargetSelection                        string                                  `json:"target_selection,omitempty"`
	AllocationMethod                       string                                  `json:"allocation_method,omitempty"`
	OncePerCustomer                        bool                                    `json:"once_per_customer,omitempty"`
	UsageLimit                             *int                                    `json:"usage_limit,omitempty"`
	StartsAt                               *time.Time                              `json:"starts_at,omitempty"`
	EndsAt                                 *time.Time                              `json:"ends_at,omitempty"`
	CreatedAt                              *time.Time                              `json:"created_at,omitempty"`
	UpdatedAt                              *time.Time                              `json:"updated_at,omitempty"`
	EntitledProductIDs                     []int64                                 `json:"entitled_product_ids,omitempty"`
	EntitledVariantIDs                     []int64                                 `json:"entitled_variant_ids,omitempty"`
	EntitledCollectionIDs                  []int64                                 `json:"entitled_collection_ids,omitempty"`
	EntitledCountryIDs                     []int64                                 `json:"entitled_country_ids,omitempty"`
	PrerequisiteProductIDs                 []int64                                 `json:"prerequisite_product_ids,omitempty"`
	PrerequisiteVariantIDs                 []int64                                 `json:"prerequisite_variant_ids,omitempty"`
	PrerequisiteCollectionIDs              []int64                                 `json:"prerequisite_collection_ids,omitempty"`
	PrerequisiteSavedSearchIDs             []int64                                 `json:"prerequisite_saved_search_ids,omitempty"`
	PrerequisiteCustomerIDs                []int64                                 `json:"prerequisite_customer_ids,omitempty"`
	PrerequisiteSubtotalRange              *PrerequisiteSubtotalRange              `json:"prerequisite_subtotal_range,omitempty"`
	PrerequisiteQuantityRange              *PrerequisiteQuantityRange              `json:"prerequisite_quantity_range,omitempty"`
	PrerequisiteShippingPriceRange         *PrerequisiteShippingPriceRange         `json:"prerequisite_shipping_price_range,omitempty"`
	PrerequisiteToEntitlementQuantityRatio *PrerequisiteToEntitlementQuantityRatio `json:"prerequisite_to_entitlement_quantity_ratio,omitempty"`
	Title                                  string                                  `json:"title,omitempty"`
}

type PrerequisiteQuantityRange struct {
	GreaterThanOrEqualTo int `json:"greater_than_or_equal_to,omitempty"`
}

type PrerequisiteShippingPriceRange struct {
	LessThanOrEqualTo float64 `json:"less_than_or_equal_to,omitempty"`
}

type PrerequisiteSubtotalRange struct {
	GreaterThanOrEqualTo float64 `json:"greater_than_or_equal_to,omitempty"`
}

type PrerequisiteToEntitlementQuantityRatio struct {
	PrerequisiteQuantity *int `json:"prerequisite_quantity,omitempty"`
	EntitledQuantity     *int `json:"entitled_quantity,omitempty"`
}

// Represents the result from the price_rules/X.json endpoint
type PriceRuleResource struct {
	PriceRule *PriceRule `json:"price_rule"`
}

// Represents the result from the price_rules.json endpoint
type PriceRulesResource struct {
	PriceRules []PriceRule `json:"price_rules"`
}

// List price rules
func (s *PriceRuleServiceOp) List(options interface{}) ([]PriceRule, error) {
	path := fmt.Sprintf("%s.json", priceRuleBasePath)
	resource := new(PriceRulesResource)
	err := s.client.Get(path, resource, options)
	return resource.PriceRules, err
}

// Get price rule
func (s *PriceRuleServiceOp) Get(priceRuleID int64, options interface{}) (*PriceRule, error) {
	path := fmt.Sprintf("%s/%d.json", priceRuleBasePath, priceRuleID)
	resource := new(PriceRuleResource)
	err := s.client.Get(path, resource, options)
	return resource.PriceRule, err
}

// Create a new price rule
func (s *PriceRuleServiceOp) Create(priceRule PriceRule) (*PriceRule, error) {
	path := fmt.Sprintf("%s.json", priceRuleBasePath)
	wrappedData := PriceRuleResource{PriceRule: &priceRule}
	resource := new(PriceRuleResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.PriceRule, err
}

// Update an existing price rule
func (s *PriceRuleServiceOp) Update(priceRule PriceRule) (*PriceRule, error) {
	path := fmt.Sprintf("%s/%d.json", priceRuleBasePath, priceRule.ID)
	wrappedData := PriceRuleResource{PriceRule: &priceRule}
	resource := new(PriceRuleResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.PriceRule, err
}

// Delete an existing price rule
func (s *PriceRuleServiceOp) Delete(priceRuleID int64) error {
	path := fmt.Sprintf("%s/%d.json", priceRuleBasePath, priceRuleID)
	return s.client.Delete(path)
}
