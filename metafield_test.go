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

func MetafieldTests(t *testing.T, metafield Metafield) {
	// Check that ID is assigned to the returned metafield
	expectedInt := int64(1)
	if metafield.ID != expectedInt {
		t.Errorf("Metafield.ID returned %+v, expected %+v", metafield.ID, expectedInt)
	}
}

func TestMetafieldList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Metafield.List(nil)
	if err != nil {
		t.Errorf("Metafield.List returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Metafield.List returned %+v, expected %+v", metafields, expected)
	}
}

func TestMetafieldListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	metafields, err := client.Metafield.List(nil)
	if metafields != nil {
		t.Errorf("Metafield.List returned metafields, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Metafield.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestMetafieldListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedMetafields []Metafield
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"metafields": [{"id":1},{"id":2}]}`,
			"",
			[]Metafield{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Metafield(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Metafield(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Metafield(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Metafield(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Metafield(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"metafields": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Metafield{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"metafields": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Metafield{{ID: 2}},
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

		metafields, pagination, err := client.Metafield.ListWithPagination(nil)
		if !reflect.DeepEqual(metafields, c.expectedMetafields) {
			t.Errorf("test %d Metafield.ListWithPagination metafields returned %+v, expected %+v", i, metafields, c.expectedMetafields)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Metafield.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Metafield.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestMetafieldCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Metafield.Count(nil)
	if err != nil {
		t.Errorf("Metafield.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Metafield.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Metafield.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Metafield.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Metafield.Count returned %d, expected %d", cnt, expected)
	}
}

func TestMetafieldGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield, err := client.Metafield.Get(1, nil)
	if err != nil {
		t.Errorf("Metafield.Get returned error: %v", err)
	}

	createdAt := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2016, time.January, 2, 0, 0, 0, 0, time.UTC)
	expected := &Metafield{
		ID:                1,
		Key:               "app_key",
		Value:             "app_value",
		ValueType:         "string",
		Namespace:         "affiliates",
		Description:       "some amaaazing app's value",
		OwnerId:           1,
		CreatedAt:         &createdAt,
		UpdatedAt:         &updatedAt,
		OwnerResource:     "shop",
		AdminGraphqlAPIID: "gid://shopify/Metafield/1",
	}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Metafield.Get returned %+v, expected %+v", metafield, expected)
	}
}

func TestMetafieldCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Namespace: "inventory",
		Key:       "warehouse",
		Value:     "25",
		ValueType: "integer",
	}

	returnedMetafield, err := client.Metafield.Create(metafield)
	if err != nil {
		t.Errorf("Metafield.Create returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestMetafieldUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        1,
		Value:     "something new",
		ValueType: "string",
	}

	returnedMetafield, err := client.Metafield.Update(metafield)
	if err != nil {
		t.Errorf("Metafield.Update returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestMetafieldDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/metafields/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Metafield.Delete(1)
	if err != nil {
		t.Errorf("Metafield.Delete returned error: %v", err)
	}
}
