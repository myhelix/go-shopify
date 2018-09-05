package goshopify

import (
	"reflect"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func DiscountCodeTests(t *testing.T, discountCode DiscountCode) {
	// Check that ID is assigned to the returned discount code
	expectedInt := int64(507328175)
	if discountCode.ID != expectedInt {
		t.Errorf("DiscountCode.ID returned %+v, expected %+v", discountCode.ID, expectedInt)
	}
}

func TestDiscountCodeList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/price_rules/1/discount_codes.json",
		httpmock.NewStringResponder(200, `{"discount_codes": [{"id":2}]}`))

	discountCodes, err := client.DiscountCode.List(1, nil)
	if err != nil {
		t.Errorf("DiscountCode.List returned error: %v", err)
	}

	expected := []DiscountCode{{ID: 2}}
	if !reflect.DeepEqual(discountCodes, expected) {
		t.Errorf("DiscountCode.List returned %+v, expected %+v", discountCodes, expected)
	}
}

func TestDiscountCodeGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/price_rules/1/discount_codes/2.json",
		httpmock.NewStringResponder(200, `{"discount_code": {"id":2}}`))

	discountCode, err := client.DiscountCode.Get(1, 2, nil)
	if err != nil {
		t.Errorf("DiscountCode.Get returned error: %v", err)
	}

	expected := &DiscountCode{ID: 2}
	if !reflect.DeepEqual(discountCode, expected) {
		t.Errorf("DiscountCode.Get returned %+v, expected %+v", discountCode, expected)
	}
}

func TestDiscountCodeCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/price_rules/1/discount_codes.json",
		httpmock.NewBytesResponder(200, loadFixture("discountcode.json")))

	discountCode := DiscountCode{
		ID:          507328175,
		PriceRuleID: 507328176,
		Code:        "SUMMERSALE10OFF",
		UsageCount:  0,
	}

	returnedDiscountCode, err := client.DiscountCode.Create(1, discountCode)
	if err != nil {
		t.Errorf("DiscountCode.Create returned error: %v", err)
	}

	DiscountCodeTests(t, *returnedDiscountCode)
}

func TestDiscountCodeUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", "https://fooshop.myshopify.com/admin/price_rules/1/discount_codes/2.json",
		httpmock.NewBytesResponder(200, loadFixture("discountcode.json")))

	discountCode := DiscountCode{
		ID:          507328175,
		PriceRuleID: 507328176,
		Code:        "SUMMERSALE10OFF",
		UsageCount:  0,
	}

	returnedDiscountCode, err := client.DiscountCode.Update(1, 2, discountCode)
	if err != nil {
		t.Errorf("DiscountCode.Update returned error: %v", err)
	}

	DiscountCodeTests(t, *returnedDiscountCode)
}

func TestDiscountCodeDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", "https://fooshop.myshopify.com/admin/price_rules/1/discount_codes/2.json",
		httpmock.NewStringResponder(200, "{}"))

	err := client.DiscountCode.Delete(1, 2)
	if err != nil {
		t.Errorf("DiscountCode.Delete returned error: %v", err)
	}
}
