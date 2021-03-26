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

func pageTests(t *testing.T, page Page) {
	// Check that ID is assigned to the returned page
	expectedInt := int64(1)
	if page.ID != expectedInt {
		t.Errorf("Page.ID returned %+v, expected %+v", page.ID, expectedInt)
	}
}

func TestPageList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"pages": [{"id":1},{"id":2}]}`))

	pages, err := client.Page.List(nil)
	if err != nil {
		t.Errorf("Page.List returned error: %v", err)
	}

	expected := []Page{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(pages, expected) {
		t.Errorf("Page.List returned %+v, expected %+v", pages, expected)
	}
}

func TestPageListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	pages, err := client.Page.List(nil)
	if pages != nil {
		t.Errorf("Page.List returned pages, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Page.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestPageListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/pages.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedPages      []Page
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"pages": [{"id":1},{"id":2}]}`,
			"",
			[]Page{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Page(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Page(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Page(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Page(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Page(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"pages": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Page{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"pages": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Page{{ID: 2}},
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

		pages, pagination, err := client.Page.ListWithPagination(nil)
		if !reflect.DeepEqual(pages, c.expectedPages) {
			t.Errorf("test %d Page.ListWithPagination pages returned %+v, expected %+v", i, pages, c.expectedPages)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Page.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Page.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestPageCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Page.Count(nil)
	if err != nil {
		t.Errorf("Page.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Page.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Page.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Page.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Page.Count returned %d, expected %d", cnt, expected)
	}
}

func TestPageGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"page": {"id":1}}`))

	page, err := client.Page.Get(1, nil)
	if err != nil {
		t.Errorf("Page.Get returned error: %v", err)
	}

	expected := &Page{ID: 1}
	if !reflect.DeepEqual(page, expected) {
		t.Errorf("Page.Get returned %+v, expected %+v", page, expected)
	}
}

func TestPageCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("page.json")))

	page := Page{
		Title:    "404",
		BodyHTML: "<strong>NOT FOUND!<\\/strong>",
	}

	returnedPage, err := client.Page.Create(page)
	if err != nil {
		t.Errorf("Page.Create returned error: %v", err)
	}

	pageTests(t, *returnedPage)
}

func TestPageUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("page.json")))

	page := Page{
		ID: 1,
	}

	returnedPage, err := client.Page.Update(page)
	if err != nil {
		t.Errorf("Page.Update returned error: %v", err)
	}

	pageTests(t, *returnedPage)
}

func TestPageDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Page.Delete(1)
	if err != nil {
		t.Errorf("Page.Delete returned error: %v", err)
	}
}

func TestPageListMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafields": [{"id":1},{"id":2}]}`))

	metafields, err := client.Page.ListMetafields(1, nil)
	if err != nil {
		t.Errorf("Page.ListMetafields() returned error: %v", err)
	}

	expected := []Metafield{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(metafields, expected) {
		t.Errorf("Page.ListMetafields() returned %+v, expected %+v", metafields, expected)
	}
}

func TestPageCountMetafields(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Page.CountMetafields(1, nil)
	if err != nil {
		t.Errorf("Page.CountMetafields() returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Page.CountMetafields() returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Page.CountMetafields(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Page.CountMetafields() returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Page.CountMetafields() returned %d, expected %d", cnt, expected)
	}
}

func TestPageGetMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"metafield": {"id":2}}`))

	metafield, err := client.Page.GetMetafield(1, 2, nil)
	if err != nil {
		t.Errorf("Page.GetMetafield() returned error: %v", err)
	}

	expected := &Metafield{ID: 2}
	if !reflect.DeepEqual(metafield, expected) {
		t.Errorf("Page.GetMetafield() returned %+v, expected %+v", metafield, expected)
	}
}

func TestPageCreateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Page.CreateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Page.CreateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestPageUpdateMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields/2.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("metafield.json")))

	metafield := Metafield{
		ID:        2,
		Key:       "app_key",
		Value:     "app_value",
		ValueType: "string",
		Namespace: "affiliates",
	}

	returnedMetafield, err := client.Page.UpdateMetafield(1, metafield)
	if err != nil {
		t.Errorf("Page.UpdateMetafield() returned error: %v", err)
	}

	MetafieldTests(t, *returnedMetafield)
}

func TestPageDeleteMetafield(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/pages/1/metafields/2.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Page.DeleteMetafield(1, 2)
	if err != nil {
		t.Errorf("Page.DeleteMetafield() returned error: %v", err)
	}
}
