package goshopify

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const giftCardsBasePath = "admin/gift_cards"

// GiftCardService is an interface for interfacing with the gift card endpoints
// of the Shopify API.
// https://help.shopify.com/en/api/reference/plus/giftcard
type GiftCardService interface {
	List(interface{}) ([]GiftCard, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*GiftCard, error)
	Search(interface{}) ([]GiftCard, error)
	Create(GiftCard) (*GiftCard, error)
	Update(GiftCard) (*GiftCard, error)
	Disable(int64) (*GiftCard, error)
}

// GiftCardServiceOp handles communication with the gift card related methods of
// the Shopify API.
type GiftCardServiceOp struct {
	client *Client
}

// GiftCard represents a Shopify gift card.
type GiftCard struct {
	ID             int64            `json:"id,omitempty"`
	InitialValue   *decimal.Decimal `json:"initial_value,omitempty"`
	Balance        *decimal.Decimal `json:"balance,omitempty"`
	Code           string           `json:"code,omitempty"`
	MaskedCode     string           `json:"masked_code,omitempty"`
	Currency       string           `json:"currency,omitempty"`
	Note           string           `json:"note,omitempty"`
	TemplateSuffix string           `json:"template_suffix,omitempty"`
	LastCharacters string           `json:"last_characters,omitempty"`
	ExpiresOn      string           `json:"expires_on,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	UpdatedAt      *time.Time       `json:"updated_at,omitempty"`
	DisabledAt     *time.Time       `json:"disabled_at,omitempty"`
	APIClientID    int64            `json:"api_client_id,omitempty"`
	OrderID        int64            `json:"order_id,omitempty"`
	UserID         int64            `json:"user_id,omitempty"`
	CustomerID     int64            `json:"customer_id,omitempty"`
	LineItemID     int64            `json:"line_item_id,omitempty"`
}

// Represents the result from the gift_card/X.json endpoint
type GiftCardResource struct {
	GiftCard *GiftCard `json:"gift_card"`
}

// Represents the result from the gift_cards/X.json endpoint
type GiftCardsResource struct {
	GiftCards []GiftCard `json:"gift_cards"`
}

// Represents the options available when searching for a gift card
type GiftCardSearchOptions struct {
	Page   int    `url:"page,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Fields string `url:"fields,omitempty"`
	Order  string `url:"order,omitempty"`
	Query  string `url:"query,omitempty"`
}

// List gift cards
func (s *GiftCardServiceOp) List(options interface{}) ([]GiftCard, error) {
	path := fmt.Sprintf("%s.json", giftCardsBasePath)
	resource := new(GiftCardsResource)
	err := s.client.Get(path, resource, options)
	return resource.GiftCards, err
}

// Count gift cards
func (s *GiftCardServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/count.json", giftCardsBasePath)
	return s.client.Count(path, options)
}

// Get gift card
func (s *GiftCardServiceOp) Get(giftCardID int64, options interface{}) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%v.json", giftCardsBasePath, giftCardID)
	resource := new(GiftCardResource)
	err := s.client.Get(path, resource, options)
	return resource.GiftCard, err
}

// Search gift cards
func (s *GiftCardServiceOp) Search(options interface{}) ([]GiftCard, error) {
	path := fmt.Sprintf("%s/search.json", giftCardsBasePath)
	resource := new(GiftCardsResource)
	err := s.client.Get(path, resource, options)
	return resource.GiftCards, err
}

// Create gift card
func (s *GiftCardServiceOp) Create(giftCard GiftCard) (*GiftCard, error) {
	path := fmt.Sprintf("%s.json", giftCardsBasePath)
	wrappedData := GiftCardResource{GiftCard: &giftCard}
	resource := new(GiftCardResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.GiftCard, err
}

// Update gift card
func (s *GiftCardServiceOp) Update(giftCard GiftCard) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%d.json", giftCardsBasePath, giftCard.ID)
	wrappedData := GiftCardResource{GiftCard: &giftCard}
	resource := new(GiftCardResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.GiftCard, err
}

// Disable gift card
func (s *GiftCardServiceOp) Disable(giftCardID int64) (*GiftCard, error) {
	path := fmt.Sprintf("%s/%d/disable.json", giftCardsBasePath, giftCardID)
	resource := new(GiftCardResource)
	err := s.client.Post(path, nil, resource)
	return resource.GiftCard, err
}
