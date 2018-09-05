package goshopify

import (
	"reflect"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func priceRuleTests(t *testing.T, priceRule PriceRule) {
	// Check that ID is assigned to the returned price rule
	expectedInt := int64(996341478)
	if priceRule.ID != expectedInt {
		t.Errorf("PriceRule.ID returned %+v, expected %+v", priceRule.ID, expectedInt)
	}
}

func TestPriceRuleList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/price_rules.json",
		httpmock.NewStringResponder(200, `{"price_rules": [{"id":1}]}`))

	priceRules, err := client.PriceRule.List(nil)
	if err != nil {
		t.Errorf("PriceRule.List returned error: %v", err)
	}

	expected := []PriceRule{{ID: 1}}
	if !reflect.DeepEqual(priceRules, expected) {
		t.Errorf("PriceRule.List returned %+v, expected %+v", priceRules, expected)
	}
}

func TestPriceRuleGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/price_rules/1.json",
		httpmock.NewStringResponder(200, `{"price_rule": {"id":1}}`))

	price_rule, err := client.PriceRule.Get(1, nil)
	if err != nil {
		t.Errorf("PriceRule.Get returned error: %v", err)
	}

	expected := &PriceRule{ID: 1}
	if !reflect.DeepEqual(price_rule, expected) {
		t.Errorf("PriceRule.Get returned %+v, expected %+v", price_rule, expected)
	}
}

func TestPriceRuleCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/price_rules.json",
		httpmock.NewBytesResponder(200, loadFixture("pricerule.json")))

	loc := time.FixedZone("AEST", 10)
	createdAt := time.Date(2017, time.September, 23, 18, 15, 47, 0, loc)
	updatedAt := time.Date(2017, time.September, 23, 18, 15, 47, 0, loc)
	startsAt := time.Date(2017, time.September, 23, 18, 15, 47, 0, loc)

	prerequisiteToEntitlementQuantityRatio := &PrerequisiteToEntitlementQuantityRatio{
		PrerequisiteQuantity: nil,
		EntitledQuantity:     nil,
	}

	priceRule := PriceRule{
		ID:                                     996341478,
		ValueType:                              "fixed_amount",
		Value:                                  "-10.0",
		CustomerSelection:                      "all",
		TargetType:                             "line_item",
		TargetSelection:                        "all",
		AllocationMethod:                       "across",
		OncePerCustomer:                        false,
		UsageLimit:                             nil,
		StartsAt:                               &startsAt,
		EndsAt:                                 nil,
		CreatedAt:                              &createdAt,
		UpdatedAt:                              &updatedAt,
		EntitledProductIDs:                     []int64{},
		EntitledVariantIDs:                     []int64{},
		EntitledCollectionIDs:                  []int64{},
		EntitledCountryIDs:                     []int64{},
		PrerequisiteProductIDs:                 []int64{},
		PrerequisiteVariantIDs:                 []int64{},
		PrerequisiteCollectionIDs:              []int64{},
		PrerequisiteSavedSearchIDs:             []int64{},
		PrerequisiteCustomerIDs:                []int64{},
		PrerequisiteSubtotalRange:              nil,
		PrerequisiteQuantityRange:              nil,
		PrerequisiteShippingPriceRange:         nil,
		PrerequisiteToEntitlementQuantityRatio: prerequisiteToEntitlementQuantityRatio,
		Title: "SUMMERSALE10OFF",
	}

	returnedPriceRule, err := client.PriceRule.Create(priceRule)
	if err != nil {
		t.Errorf("PriceRule.Create returned error: %v", err)
	}

	priceRuleTests(t, *returnedPriceRule)
}

func TestPriceRuleUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", "https://fooshop.myshopify.com/admin/price_rules/1.json",
		httpmock.NewBytesResponder(200, loadFixture("pricerule.json")))

	priceRule := PriceRule{
		ID: 1,
	}

	returnedPriceRule, err := client.PriceRule.Update(priceRule)
	if err != nil {
		t.Errorf("PriceRule.Update returned error: %v", err)
	}

	priceRuleTests(t, *returnedPriceRule)
}

func TestPriceRuleDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", "https://fooshop.myshopify.com/admin/price_rules/1.json",
		httpmock.NewStringResponder(200, "{}"))

	err := client.PriceRule.Delete(1)
	if err != nil {
		t.Errorf("PriceRule.Delete returned error: %v", err)
	}
}
