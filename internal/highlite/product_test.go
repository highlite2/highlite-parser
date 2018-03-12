package highlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImages(t *testing.T) {
	// arrange
	testData := "40462.jpg|40462_ .jpg|40462_detail1.jpg|40462_detail2.jpg|40462_detail3.jpg|40462_detail4.jpg|40462_detail5.jpg"
	expected := []string{"40462.jpg", "40462_ .jpg", "40462_detail1.jpg", "40462_detail2.jpg", "40462_detail3.jpg", "40462_detail4.jpg", "40462_detail5.jpg"}

	// act
	actual := parseImages(testData)

	// assert
	assert.Equal(t, expected, actual)
}

func TestParseImages_Empty(t *testing.T) {
	// arrange
	testData := ""
	expected := []string{}

	// act
	actual := parseImages(testData)

	// assert
	assert.Equal(t, expected, actual)
}

func TestProduct_ProductName(t *testing.T) {
	// arrange
	product1 := Product{Name: "name"}
	product2 := Product{Name: "name", Brand: "brand"}

	// assert
	assert.Equal(t, "name", product1.ProductName())
	assert.Equal(t, "brand name", product2.ProductName())
}
