package goshopify

import (
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func imageTests(t *testing.T, image Image) {
	// Check that ID is set
	expectedImageID := int64(1)
	if image.ID != expectedImageID {
		t.Errorf("Image.ID returned %+v, expected %+v", image.ID, expectedImageID)
	}

	// Check that product_id is set
	expectedProductID := int64(1)
	if image.ProductID != expectedProductID {
		t.Errorf("Image.ProductID returned %+v, expected %+v", image.ProductID, expectedProductID)
	}

	// Check that position is set
	expectedPosition := 1
	if image.Position != expectedPosition {
		t.Errorf("Image.Position returned %+v, expected %+v", image.Position, expectedPosition)
	}

	// Check that width is set
	expectedWidth := 123
	if image.Width != expectedWidth {
		t.Errorf("Image.Width returned %+v, expected %+v", image.Width, expectedWidth)
	}

	// Check that height is set
	expectedHeight := 456
	if image.Height != expectedHeight {
		t.Errorf("Image.Height returned %+v, expected %+v", image.Height, expectedHeight)
	}

	// Check that src is set
	expectedSrc := "https://cdn.shopify.com/s/files/1/0006/9093/3842/products/ipod-nano.png?v=1500937783"
	if image.Src != expectedSrc {
		t.Errorf("Image.Src returned %+v, expected %+v", image.Src, expectedSrc)
	}

	// Check that variant ids are set
	expectedVariantIDs := make([]int64, 2)
	expectedVariantIDs[0] = int64(808950810)
	expectedVariantIDs[1] = int64(808950811)

	if image.VariantIDs[0] != expectedVariantIDs[0] {
		t.Errorf("Image.VariantIDs[0] returned %+v, expected %+v", image.VariantIDs[0], expectedVariantIDs[0])
	}
	if image.VariantIDs[1] != expectedVariantIDs[1] {
		t.Errorf("Image.VariantIDs[0] returned %+v, expected %+v", image.VariantIDs[1], expectedVariantIDs[1])
	}

	// Check that CreatedAt date is set
	expectedCreatedAt := time.Date(2017, time.July, 24, 19, 9, 43, 0, time.UTC)
	if !expectedCreatedAt.Equal(*image.CreatedAt) {
		t.Errorf("Image.CreatedAt returned %+v, expected %+v", image.CreatedAt, expectedCreatedAt)
	}

	// Check that UpdatedAt date is set
	expectedUpdatedAt := time.Date(2017, time.July, 24, 19, 9, 43, 0, time.UTC)
	if !expectedUpdatedAt.Equal(*image.UpdatedAt) {
		t.Errorf("Image.UpdatedAt returned %+v, expected %+v", image.UpdatedAt, expectedUpdatedAt)
	}
}

func TestImageList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/products/1/images.json",
		httpmock.NewBytesResponder(200, loadFixture("images.json")))

	images, err := client.Image.List(1, nil)
	if err != nil {
		t.Errorf("Images.List returned error: %v", err)
	}

	// Check that images were parsed
	if len(images) != 2 {
		t.Errorf("Image.List got %v images, expected 2", len(images))
	}

	imageTests(t, images[0])
}

func TestImageCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/products/1/images/count.json",
		httpmock.NewStringResponder(200, `{"count": 2}`))

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/products/1/images/count.json?created_at_min=2016-01-01T00%3A00%3A00Z",
		httpmock.NewStringResponder(200, `{"count": 1}`))

	cnt, err := client.Image.Count(1, nil)
	if err != nil {
		t.Errorf("Image.Count returned error: %v", err)
	}

	expected := 2
	if cnt != expected {
		t.Errorf("Image.Count returned %d, expected %d", cnt, expected)
	}

	date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	cnt, err = client.Image.Count(1, CountOptions{CreatedAtMin: date})
	if err != nil {
		t.Errorf("Image.Count returned %d, expected %d", cnt, expected)
	}

	expected = 1
	if cnt != expected {
		t.Errorf("Image.Count returned %d, expected %d", cnt, expected)
	}
}

func TestImageGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://fooshop.myshopify.com/admin/products/1/images/1.json",
		httpmock.NewBytesResponder(200, loadFixture("image.json")))

	image, err := client.Image.Get(1, 1, nil)
	if err != nil {
		t.Errorf("Image.Get returned error: %v", err)
	}

	imageTests(t, *image)
}

func TestImageCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/products/1/images.json",
		httpmock.NewBytesResponder(200, loadFixture("image.json")))

	variantIDs := make([]int64, 2)
	variantIDs[0] = int64(808950810)
	variantIDs[1] = int64(808950811)

	image := Image{
		Src:        "https://cdn.shopify.com/s/files/1/0006/9093/3842/products/ipod-nano.png?v=1500937783",
		VariantIDs: variantIDs,
	}
	returnedImage, err := client.Image.Create(1, image)
	if err != nil {
		t.Errorf("Image.Create returned error %v", err)
	}

	imageTests(t, *returnedImage)
}

func TestImageUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", "https://fooshop.myshopify.com/admin/products/1/images/1.json",
		httpmock.NewBytesResponder(200, loadFixture("image.json")))

	// Take an existing image
	variantIDs := make([]int64, 2)
	variantIDs[0] = int64(808950810)
	variantIDs[1] = int64(457924702)
	existingImage := Image{
		ID:         1,
		VariantIDs: variantIDs,
	}
	// And update it
	existingImage.VariantIDs[1] = 808950811
	returnedImage, err := client.Image.Update(1, existingImage)
	if err != nil {
		t.Errorf("Image.Update returned error %v", err)
	}

	imageTests(t, *returnedImage)
}

func TestImageDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", "https://fooshop.myshopify.com/admin/products/1/images/1.json",
		httpmock.NewStringResponder(200, "{}"))

	err := client.Image.Delete(1, 1)
	if err != nil {
		t.Errorf("Image.Delete returned error: %v", err)
	}
}
