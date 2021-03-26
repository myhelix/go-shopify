package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGiftCardList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"gift_cards": [{"id":1},{"id":2}]}`,
		))

	giftCards, err := client.GiftCard.List(nil)
	if err != nil {
		t.Errorf("GiftCard.List returned error: %v", err)
	}

	expected := []GiftCard{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(giftCards, expected) {
		t.Errorf("GiftCard.List returned %+v, expected %+v", giftCards, expected)
	}
}

func TestGiftCardCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/count.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"count": 5}`,
		))

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/count.json?created_at_min=2016-01-01T00%%3A00%%3A00Z", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"count": 2}`,
		))

	cnt, err := client.GiftCard.Count(nil)
	if err != nil {
		t.Errorf("GiftCard.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("GiftCard.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.GiftCard.Count(CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("GiftCard.Count returned error: %v", err)
	}

	expected = 2
	if cnt != expected {
		t.Errorf("GiftCard.Count returned %d, expected %d", cnt, expected)
	}
}

func TestGiftCardGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"gift_card": {"id":1}}`,
		))

	giftCard, err := client.GiftCard.Get(1, nil)
	if err != nil {
		t.Errorf("GiftCard.Get returned error: %v", err)
	}

	expected := &GiftCard{ID: 1}
	if !reflect.DeepEqual(giftCard, expected) {
		t.Errorf("GiftCard.Get returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardSearch(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/search.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"gift_cards": [{"id":1},{"id":2}]}`,
		))

	giftCards, err := client.GiftCard.Search(nil)
	if err != nil {
		t.Errorf("GiftCard.Search returned error: %v", err)
	}

	expected := []GiftCard{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(giftCards, expected) {
		t.Errorf("GiftCard.Search returned %+v, expected %+v", giftCards, expected)
	}
}

func TestGiftCardCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card.json"),
		))

	giftCard := GiftCard{}
	returnedGiftCard, err := client.GiftCard.Create(giftCard)
	if err != nil {
		t.Errorf("GiftCard.Create returned error: %v", err)
	}

	expectedCustomerID := int64(1)
	if returnedGiftCard.ID != expectedCustomerID {
		t.Errorf("GiftCard.InitialValue returned %+v expected %+v", returnedGiftCard.ID, expectedCustomerID)
	}
}

func TestGiftCardUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card.json"),
		))

	giftCard := GiftCard{
		ID: 1,
	}

	returnedGiftCard, err := client.GiftCard.Update(giftCard)
	if err != nil {
		t.Errorf("GiftCard.Update returned error: %v", err)
	}

	expectedCustomerID := int64(1)
	if returnedGiftCard.ID != expectedCustomerID {
		t.Errorf("GiftCard.InitialValue returned %+v expected %+v", returnedGiftCard.ID, expectedCustomerID)
	}
}

func TestGiftCardDisable(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1/disable.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card.json"),
		))

	giftCardID := int64(1)
	returnedGiftCard, err := client.GiftCard.Disable(giftCardID)
	if err != nil {
		t.Errorf("GiftCard.Disable returned error: %v", err)
	}

	expectedCustomerID := int64(1)
	if returnedGiftCard.ID != expectedCustomerID {
		t.Errorf("GiftCard.InitialValue returned %+v expected %+v", returnedGiftCard.ID, expectedCustomerID)
	}
}
