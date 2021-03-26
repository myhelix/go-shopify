package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestPriceRuleGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Get(1)
	if err != nil {
		t.Errorf("PriceRule.Get returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Get returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/list.json"),
		),
	)

	rules, err := client.PriceRule.List()
	if err != nil {
		t.Errorf("PriceRule.List returned error: %v", err)
	}

	expected := []PriceRule{{ID: 1}}
	if expected[0].ID != rules[0].ID {
		t.Errorf("PriceRule.List returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	priceRules, err := client.PriceRule.List()
	if priceRules != nil {
		t.Errorf("PriceRule.List returned price rules, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("PriceRule.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestPriceRuleListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedPriceRules   []PriceRule
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"priceRules": [{"id":1},{"id":2}]}`,
			"",
			[]PriceRule{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]PriceRule(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]PriceRule(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]PriceRule(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]PriceRule(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]PriceRule(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"priceRules": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]PriceRule{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"priceRules": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]PriceRule{{ID: 2}},
			&Pagination{
				NextPageOptions:     &ListOptions{PageInfo: "foo"},
				PreviousPageOptions: &ListOptions{PageInfo: "bar"},
			},
			nil,
		},
	}
	for i, c := range cases {
		response := &http.Response{
			StatusCode: 200,
			Body:       httpmock.NewRespBodyFromString(c.body),
			Header: http.Header{
				"Link": {c.linkHeader},
			},
		}

		httpmock.RegisterResponder("GET", listURL, httpmock.ResponderFromResponse(response))

		priceRules, pagination, err := client.PriceRule.ListWithPagination(nil)
		if !reflect.DeepEqual(priceRules, c.expectedPriceRules) {
			t.Errorf("test %d PriceRule.ListWithPagination price rules returned %+v, expected %+v", i, priceRules, c.expectedPriceRules)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d PriceRule.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d PriceRule.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestPriceRuleCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Create(PriceRule{})
	if err != nil {
		t.Errorf("PriceRule.Create returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Create returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	rules, err := client.PriceRule.Update(PriceRule{ID: 1})
	if err != nil {
		t.Errorf("PriceRule.Update returned error: %v", err)
	}

	expected := PriceRule{ID: 1}
	if expected.ID != rules.ID {
		t.Errorf("PriceRule.Update returned %+v, expected %+v", rules, expected)
	}
}

func TestPriceRuleDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/price_rules/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("price_rule/get.json"),
		),
	)

	err := client.PriceRule.Delete(1)
	if err != nil {
		t.Errorf("PriceRule.Delete returned error: %v", err)
	}
}

func TestPriceRuleSetters(t *testing.T) {
	pr := PriceRule{}
	prereqSubtotalRange := "1.5"
	prereqQuantityRange := 2
	prereqShippingPrice := "5.5"
	prereqRatioQuantity := 1
	prereqRatioEntitledQuantity := 1
	badMoneyString := "dog"

	// Test bad money strings
	err := pr.SetPrerequisiteSubtotalRange(&badMoneyString)
	if err == nil {
		t.Errorf("Expected error from setting bad string as prerequisite subtotal range: %s", badMoneyString)
	}

	err = pr.SetPrerequisiteShippingPriceRange(&badMoneyString)
	if err == nil {
		t.Errorf("Expected error from setting bad string as prerequisite shipping price: %s", badMoneyString)
	}

	// Test populating values
	err = pr.SetPrerequisiteSubtotalRange(&prereqSubtotalRange)
	if err != nil {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteQuantityRange(&prereqQuantityRange)
	err = pr.SetPrerequisiteShippingPriceRange(&prereqShippingPrice)
	if err != nil {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(&prereqRatioQuantity, &prereqRatioEntitledQuantity)

	if pr.PrerequisiteSubtotalRange.GreaterThanOrEqualTo != prereqSubtotalRange {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	if pr.PrerequisiteQuantityRange.GreaterThanOrEqualTo != prereqQuantityRange {
		t.Errorf("Failed to set prerequisite quantity range: %d", prereqQuantityRange)
	}

	if pr.PrerequisiteShippingPriceRange.LessThanOrEqualTo != prereqShippingPrice {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqShippingPrice)
	}

	if pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != prereqRatioQuantity {
		t.Errorf("Failed to set prerequisite ratio quantity: %d", prereqRatioQuantity)
	}

	if pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != prereqRatioEntitledQuantity {
		t.Errorf("Failed to set prerequisite ratio entitled quantity: %d", prereqRatioEntitledQuantity)
	}

	// Test clearing values by setting nil
	err = pr.SetPrerequisiteSubtotalRange(nil)
	if err != nil {
		t.Errorf("Failed to set prerequisite subtotal range: %s", prereqSubtotalRange)
	}

	pr.SetPrerequisiteQuantityRange(nil)
	err = pr.SetPrerequisiteShippingPriceRange(nil)
	if err != nil {
		t.Errorf("Failed to set prerequisite shipping price: %s", prereqSubtotalRange)
	}

	if pr.PrerequisiteSubtotalRange != nil {
		t.Errorf("Failed to clear prerequisite subtotal range")
	}

	if pr.PrerequisiteQuantityRange != nil {
		t.Errorf("Failed to clear prerequisite quantity range")
	}

	if pr.PrerequisiteShippingPriceRange != nil {
		t.Errorf("Failed to clear prerequisite shipping price")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(nil, &prereqRatioEntitledQuantity)
	if pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != 0 || pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != prereqRatioEntitledQuantity {
		t.Errorf("Failed to clear prerequisite-to-entitlement-quantity-ratio prerequisite quantity")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(&prereqRatioQuantity, nil)
	if pr.PrerequisiteToEntitlementQuantityRatio.EntitledQuantity != 0 || pr.PrerequisiteToEntitlementQuantityRatio.PrerequisiteQuantity != prereqRatioQuantity {
		t.Errorf("Failed to clear prerequisite-to-entitlement-quantity-ratio entitled quantity")
	}

	pr.SetPrerequisiteToEntitlementQuantityRatio(nil, nil)
	if pr.PrerequisiteToEntitlementQuantityRatio != nil {
		t.Errorf("Failed to clear wholly prerequisite to entitlement quantity ratio")
	}
}
