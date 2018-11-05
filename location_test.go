package goshopify

import (
	"reflect"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestLocationList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/locations.json",
		httpmock.NewStringResponder(200, `{"locations": [{"id":1},{"id":2}]}`))

	locations, err := client.Location.List(nil)
	if err != nil {
		t.Errorf("Location.List returned error: %v", err)
	}

	expected := []Location{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(locations, expected) {
		t.Errorf("Location.List returned %+v, expected %+v", locations, expected)
	}
}

func TestLocationCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/locations/count.json",
		httpmock.NewStringResponder(200, `{"count": 5}`))

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/locations/count.json?created_at_min=2016-01-01T00%3A00%3A00Z",
		httpmock.NewStringResponder(200, `{"count": 2}`))

	cnt, err := client.Location.Count(nil)
	if err != nil {
		t.Errorf("Location.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("Location.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Location.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Location.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("Location.Count returned %d, expected %d", cnt, expected)
	}
}

func TestLocationGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/locations/1.json",
		httpmock.NewStringResponder(200, `{"location": {"id":1}}`))

	location, err := client.Location.Get(1, nil)
	if err != nil {
		t.Errorf("Location.Get returned error: %v", err)
	}

	expected := &Location{ID: 1}
	if !reflect.DeepEqual(location, expected) {
		t.Errorf("Location.Get returned %+v, expected %+v", location, expected)
	}
}
