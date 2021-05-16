package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product defines the datastore structure for a product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

// custom validation function
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of this format: abc-defg-hijk
	reg := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := reg.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON.
// NewEncoder provides better perf than json.Unmarshal as it does
// not have to buffer the output into an in memory slice of bytes.
// This reduces allocations and the overhead of the service
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// returns a list of products
func GetProducts() Products{
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error{
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

// generates a new id for each new record inserted to the data store
func getNextID() int {
	lp := productList[len(productList) - 1]
	return lp.ID + 1
}


// productList is a hardcoded list of products and represents the data store (i.e. database)
var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Milky coffee",
		SKU:         "qwerty123",
		Price:       2.49,
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	&Product{
		ID: 				 2,
		Name:        "Colombian",
		Description: "Dark and strong coffee",
		SKU:         "asdf4576",
		Price:       1.99,
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}