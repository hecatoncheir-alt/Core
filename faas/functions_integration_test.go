package faas

import (
	"github.com/hecatoncheir/Core/configuration"
	"github.com/hecatoncheir/Storage"
	"testing"
)

func TestFunctions_MVideoPageParser(t *testing.T) {
	config := configuration.New()

	faas := Functions{
		config.APIVersion,
		config.FunctionsGateway,
		config.DatabaseGateway}

	instruction := storage.PageInstruction{
		ItemSelector:               ".c-product-tile",
		PreviewImageOfItemSelector: ".c-product-tile-picture__link .lazy-load-image-holder img",
		NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
		LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
		PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"}

	pageForParseIRI := "https://www.mvideo.ru/smartfony-i-svyaz/smartfony-205/f/page=20"

	products := faas.MVideoPageParser(pageForParseIRI, instruction)

	if len(products) != 12 {
		t.Fatalf("Expected: 12 parsed products from page, but get: %v", len(products))
	}
}

func TestFunctions_MVideoPagesCountParser(t *testing.T) {
	config := configuration.New()

	faas := Functions{
		config.APIVersion,
		config.FunctionsGateway,
		config.DatabaseGateway}

	instruction := storage.PageInstruction{
		PageInPaginationSelector: ".c-pagination > .c-pagination__num"}

	pageForParseIRI := "https://www.mvideo.ru/noutbuki-planshety-komputery/noutbuki-118"

	pagesCount := faas.MVideoPagesCountParser(pageForParseIRI, instruction)

	expected := 87

	if pagesCount != expected {
		t.Fatalf("Expected: %v parsed pages count, but get: %v", expected, pagesCount)
	}
}
