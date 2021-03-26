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

func redirectTests(t *testing.T, redirect Redirect) {
	// Check that ID is assigned to the returned redirect
	expectedInt := int64(1)
	if redirect.ID != expectedInt {
		t.Errorf("Redirect.ID returned %+v, expected %+v", redirect.ID, expectedInt)
	}
}

func TestRedirectList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"redirects": [{"id":1},{"id":2}]}`))

	redirects, err := client.Redirect.List(nil)
	if err != nil {
		t.Errorf("Redirect.List returned error: %v", err)
	}

	expected := []Redirect{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(redirects, expected) {
		t.Errorf("Redirect.List returned %+v, expected %+v", redirects, expected)
	}
}

func TestRedirectListError(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects.json", client.pathPrefix),
		httpmock.NewStringResponder(500, ""))

	expectedErrMessage := "Unknown Error"

	redirects, err := client.Redirect.List(nil)
	if redirects != nil {
		t.Errorf("Redirect.List returned redirects, expected nil: %v", err)
	}

	if err == nil || err.Error() != expectedErrMessage {
		t.Errorf("Redirect.List err returned %+v, expected %+v", err, expectedErrMessage)
	}
}

func TestRedirectListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body               string
		linkHeader         string
		expectedRedirects  []Redirect
		expectedPagination *Pagination
		expectedErr        error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"redirects": [{"id":1},{"id":2}]}`,
			"",
			[]Redirect{{ID: 1}, {ID: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]Redirect(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]Redirect(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]Redirect(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]Redirect(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]Redirect(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"redirects": [{"id":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]Redirect{{ID: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"redirects": [{"id":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]Redirect{{ID: 2}},
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

		redirects, pagination, err := client.Redirect.ListWithPagination(nil)
		if !reflect.DeepEqual(redirects, c.expectedRedirects) {
			t.Errorf("test %d Redirect.ListWithPagination redirects returned %+v, expected %+v", i, redirects, c.expectedRedirects)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d Redirect.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d Redirect.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}

func TestRedirectCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects/count.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"count": 3}`))

	params := map[string]string{"created_at_min": "2016-01-01T00:00:00Z"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects/count.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Redirect.Count(nil)
	if err != nil {
		t.Errorf("Redirect.Count returned error: %v", err)
	}

	expected := 3
	if cnt != expected {
		t.Errorf("Redirect.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Redirect.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Redirect.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Redirect.Count returned %d, expected %d", cnt, expected)
	}
}

func TestRedirectGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, `{"redirect": {"id":1}}`))

	redirect, err := client.Redirect.Get(1, nil)
	if err != nil {
		t.Errorf("Redirect.Get returned error: %v", err)
	}

	expected := &Redirect{ID: 1}
	if !reflect.DeepEqual(redirect, expected) {
		t.Errorf("Redirect.Get returned %+v, expected %+v", redirect, expected)
	}
}

func TestRedirectCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("redirect.json")))

	redirect := Redirect{
		Path:   "/from",
		Target: "/to",
	}

	returnedRedirect, err := client.Redirect.Create(redirect)
	if err != nil {
		t.Errorf("Redirect.Create returned error: %v", err)
	}

	redirectTests(t, *returnedRedirect)
}

func TestRedirectUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("redirect.json")))

	redirect := Redirect{
		ID: 1,
	}

	returnedRedirect, err := client.Redirect.Update(redirect)
	if err != nil {
		t.Errorf("Redirect.Update returned error: %v", err)
	}

	redirectTests(t, *returnedRedirect)
}

func TestRedirectDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/redirects/1.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Redirect.Delete(1)
	if err != nil {
		t.Errorf("Redirect.Delete returned error: %v", err)
	}
}
