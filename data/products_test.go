package data

import "testing"

func TestChecksValidationWithoutValues(t *testing.T) {
	p := &Product{}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestChecksValidationWithValues(t *testing.T) {
	p := &Product{
		Name:  "Adi",
		Price: 1.00,
		SKU:   "abc-def-ghi",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestChecksValidationWithWrongSKU(t *testing.T) {
	p := &Product{
		Name:  "Adi",
		Price: 1.00,
		SKU:   "abc-defghi",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
