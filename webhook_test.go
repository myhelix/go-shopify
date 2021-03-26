package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func webhookTests(t *testing.T, webhook Webhook) {
	// Check that dates are parsed
	d := time.Date(2016, time.June, 1, 14, 10, 44, 0, time.UTC)
	if !d.Equal(*webhook.CreatedAt) {
		t.Errorf("Webhook.CreatedAt returned %+v, expected %+v", webhook.CreatedAt, d)
	}

	expectedStr := "http://apple.com"
	if webhook.Address != expectedStr {
		t.Errorf("Webhook.Address returned %+v, expected %+v", webhook.Address, expectedStr)
	}

	expectedStr = "orders/create"
	if webhook.Topic != expectedStr {
		t.Errorf("Webhook.Topic returned %+v, expected %+v", webhook.Topic, expectedStr)
	}

	expectedArr := []string{"id", "updated_at"}
	if !reflect.DeepEqual(webhook.Fields, expectedArr) {
		t.Errorf("Webhook.Fields returned %+v, expected %+v", webhook.Fields, expectedArr)
	}

	expectedArr = []string{"google", "inventory"}
	if !reflect.DeepEqual(webhook.MetafieldNamespaces, expectedArr) {
		t.Errorf("Webhook.Fields returned %+v, expected %+v", webhook.MetafieldNamespaces, expectedArr)
	}
}

func TestWebhookList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("webhooks.json")))

	webhooks, err := client.Webhook.List(nil)
	if err != nil {
		t.Errorf("Webhook.List returned error: %v", err)
	}

	// Check that webhooks were parsed
	if len(webhooks) != 1 {
		t.Errorf("Webhook.List got %v webhooks, expected: 1", len(webhooks))
	}

	webhookTests(t, webhooks[0])
}

func TestWebhookListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	webhooks, err := client.Webhook.List(nil)
	if webhooks != nil {
		t.Errorf("Webhook.List returned webhooks, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Webhook.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestWebhookListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedWebhooks   []Webhook
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"webhooks": [{"id":1},{"id":2}]}`,
			"",
			[]Webhook{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Webhook(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Webhook(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Webhook(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Webhook(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Webhook(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"webhooks": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Webhook{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"webhooks": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Webhook{{ID: 2}},
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

		webhooks, pagination, err := client.Webhook.ListWithPagination(nil)
		if !reflect.DeepEqual(webhooks, c.expectedWebhooks) {
			t.Errorf("test %d Webhook.ListWithPagination webhooks returned %+v, expected %+v", i, webhooks, c.expectedWebhooks)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Webhook.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Webhook.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestWebhookGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks/4759306.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("webhook.json")))

	webhook, err := client.Webhook.Get(4759306, nil)
	if err != nil {
		t.Errorf("Webhook.Get returned error: %v", err)
	}

	webhookTests(t, *webhook)
}

func TestWebhookCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 7}`))

	params := map[string]string{"topic": "orders/paid"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Webhook.Count(nil)
	if err != nil {
		t.Errorf("Webhook.Count returned error: %v", err)
	}

	expected := 7
	if cnt != expected {
		t.Errorf("Webhook.Count returned %d, expected %d", cnt, expected)
	}

	options := WebhookOptions{Topic: "orders/paid"}
	cnt, err = client.Webhook.Count(options)
	if err != nil {
		t.Errorf("Webhook.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Webhook.Count returned %d, expected %d", cnt, expected)
	}
}

func TestWebhookCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("webhook.json")))

	webhook := Webhook{
		Topic:   "orders/create",
		Address: "http://example.com",
	}

	returnedWebhook, err := client.Webhook.Create(webhook)
	if err != nil {
		t.Errorf("Webhook.Create returned error: %v", err)
	}

	webhookTests(t, *returnedWebhook)
}

func TestWebhookUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks/4759306.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("webhook.json")))

	webhook := Webhook{
		ID:      4759306,
		Topic:   "orders/create",
		Address: "http://example.com",
	}

	returnedWebhook, err := client.Webhook.Update(webhook)
	if err != nil {
		t.Errorf("Webhook.Update returned error: %v", err)
	}

	webhookTests(t, *returnedWebhook)
}

func TestWebhookDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/webhooks/4759306.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Webhook.Delete(4759306)
	if err != nil {
		t.Errorf("Webhook.Delete returned error: %v", err)
	}
}
