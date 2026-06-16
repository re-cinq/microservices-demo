package main

import (
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func makeProduct(id, name string, categories []string) productView {
	return productView{
		Item:  &pb.Product{Id: id, Name: name, Categories: categories},
		Price: &pb.Money{CurrencyCode: "USD", Units: 10},
	}
}

func TestGroupByCategory_BasicGrouping(t *testing.T) {
	products := []productView{
		makeProduct("1", "Sunglasses", []string{"accessories"}),
		makeProduct("2", "Watch", []string{"accessories"}),
		makeProduct("3", "Mug", []string{"kitchen"}),
	}
	groups := groupByCategory(products)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Slug != "accessories" || len(groups[0].Products) != 2 {
		t.Errorf("unexpected first group: %+v", groups[0])
	}
	if groups[1].Slug != "kitchen" || len(groups[1].Products) != 1 {
		t.Errorf("unexpected second group: %+v", groups[1])
	}
}

func TestGroupByCategory_UsesFirstCategory(t *testing.T) {
	products := []productView{
		makeProduct("1", "Tank Top", []string{"clothing", "tops"}),
	}
	groups := groupByCategory(products)
	if len(groups) != 1 || groups[0].Slug != "clothing" {
		t.Errorf("expected slug 'clothing', got %+v", groups)
	}
}

func TestGroupByCategory_EmptyCategoriesFallback(t *testing.T) {
	products := []productView{
		makeProduct("1", "Mystery Item", []string{}),
	}
	groups := groupByCategory(products)
	if len(groups) != 1 || groups[0].Slug != "other" {
		t.Errorf("expected slug 'other', got %+v", groups)
	}
}

func TestGroupByCategory_PreservesOrder(t *testing.T) {
	products := []productView{
		makeProduct("1", "Mug", []string{"kitchen"}),
		makeProduct("2", "Sunglasses", []string{"accessories"}),
		makeProduct("3", "Loafers", []string{"footwear"}),
	}
	groups := groupByCategory(products)
	slugs := make([]string, len(groups))
	for i, g := range groups {
		slugs[i] = g.Slug
	}
	expected := []string{"kitchen", "accessories", "footwear"}
	for i, s := range expected {
		if slugs[i] != s {
			t.Errorf("order mismatch at %d: want %s got %s", i, s, slugs[i])
		}
	}
}

func TestGroupByCategory_Empty(t *testing.T) {
	groups := groupByCategory(nil)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for nil input, got %d", len(groups))
	}
}
