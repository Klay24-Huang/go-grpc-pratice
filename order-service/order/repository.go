package order

import (
	"fmt"
)

type Repository interface {
	FindAll() (*[]Order, error)
	FindById(Id int) (*Order, error)
	Create(order Order) (*Order, error)
	Update(order Order) (*Order, error)
	Delete(Id int) (*Order, error)
}

type repository struct {
	monkOrders []Order
}

func NewRepository() *repository {
	monkOrders := []Order{
		{1, 1, []OrderProduct{{1, 1}}},
		{2, 2, []OrderProduct{
			{2, 2},
			{3, 10},
		}},
	}
	repository := repository{
		monkOrders,
	}
	return &repository
}

func (r *repository) FindAll() (*[]Order, error) {
	return &r.monkOrders, nil
}

func (r *repository) FindById(Id int) (*Order, error) {
	var Order Order
	err := r.findOrderFromMonk(&Order, Id)

	return &Order, err
}

func (r *repository) Create(Order Order) (*Order, error) {
	Order.Id = r.monkOrders[len(r.monkOrders)-1].Id + 1
	r.monkOrders = append(r.monkOrders, Order)
	return &Order, nil
}

func (r *repository) Update(updateOrder Order) (*Order, error) {
	for i, u := range r.monkOrders {
		if u.Id == updateOrder.Id {
			r.monkOrders[i] = updateOrder
			return &updateOrder, nil
		}
	}
	return nil, fmt.Errorf("Order not found")
}

func (r *repository) Delete(Id int) (*Order, error) {
	for i, u := range r.monkOrders {
		if u.Id == Id {
			r.monkOrders = append(r.monkOrders[:i], r.monkOrders[i+1:]...)
			return &u, nil
		}
	}
	return nil, fmt.Errorf("Order not found")
}

func (r *repository) findOrderFromMonk(Order *Order, Id int) error {
	for _, u := range r.monkOrders {
		if u.Id == Id {
			*Order = u
			return nil
		}
	}

	return fmt.Errorf("Order id:%d was not found", Id)
}
