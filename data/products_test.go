package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "latte",
		Price: 1.49,
		SKU: "abd-vgh-dgs",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}