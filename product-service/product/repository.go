package product

import (
	"fmt"
)

type Repository interface {
	FindAll() (*[]Product, error)
	FindById(Id int) (*Product, error)
	Create(product Product) (*Product, error)
	Update(product Product) (*Product, error)
	Delete(Id int) (*Product, error)
}

type repository struct {
	monkProducts []Product
}

func NewRepository() *repository {
	monkProducts := []Product{
		{1, "laptop", 200},
		{2, "smart phone", 150},
		{3, "potato", 10},
	}
	repository := repository{
		monkProducts,
	}
	return &repository
}

func (r *repository) FindAll() (*[]Product, error) {
	return &r.monkProducts, nil
}

func (r *repository) FindById(Id int) (*Product, error) {
	var Product Product
	err := r.findProductFromMonk(&Product, Id)

	return &Product, err
}

func (r *repository) Create(Product Product) (*Product, error) {
	Product.Id = r.monkProducts[len(r.monkProducts)-1].Id + 1
	r.monkProducts = append(r.monkProducts, Product)
	return &Product, nil
}

func (r *repository) Update(updateProduct Product) (*Product, error) {
	for i, u := range r.monkProducts {
		if u.Id == updateProduct.Id {
			r.monkProducts[i] = updateProduct
			return &updateProduct, nil
		}
	}
	return nil, fmt.Errorf("Product not found")
}

func (r *repository) Delete(Id int) (*Product, error) {
	for i, u := range r.monkProducts {
		if u.Id == Id {
			r.monkProducts = append(r.monkProducts[:i], r.monkProducts[i+1:]...)
			return &u, nil
		}
	}
	return nil, fmt.Errorf("Product not found")
}

func (r *repository) findProductFromMonk(Product *Product, Id int) error {
	for _, u := range r.monkProducts {
		if u.Id == Id {
			*Product = u
			return nil
		}
	}

	return fmt.Errorf("Product id:%d was not found", Id)
}
