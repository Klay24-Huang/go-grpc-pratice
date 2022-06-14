package order

import (
	"context"
	"fmt"
	productPb "grpc-practice/grpc-gateway/product"
	userPb "grpc-practice/grpc-gateway/user"
	"time"
)

type Service interface {
	FindAll() (*[]OrderResponse, error)
	FindById(Id int) (*OrderResponse, error)
	Create(orderRequest OrderRequest) (*Order, error)
	Update(Id int, orderRequest OrderRequest) (*Order, error)
	Delete(Id int) (*Order, error)
}

type service struct {
	orderRepository   Repository
	userGrpcClient    userPb.UserClient
	productGrpcClient productPb.ProductClient
}

func NewService(orderRepository Repository, userClient userPb.UserClient, productClient productPb.ProductClient) *service {
	return &service{orderRepository, userClient, productClient}
}

func (s *service) FindAll() (*[]OrderResponse, error) {
	rawOrders, err := s.orderRepository.FindAll()
	var ordersResponses []OrderResponse

	for _, o := range *rawOrders {
		orderResponse := s.getOrderDetail(&o)
		ordersResponses = append(ordersResponses, *orderResponse)
	}
	return &ordersResponses, err
}

func (s *service) FindById(Id int) (*OrderResponse, error) {
	order, err := s.orderRepository.FindById(Id)

	if err != nil {
		return nil, fmt.Errorf("Order not found")
	}

	orderResponse := s.getOrderDetail(order)
	return orderResponse, err
}

func (s *service) getOrderDetail(order *Order) *OrderResponse {
	orderResponse := OrderResponse{Id: order.Id}
	// get user info
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := s.userGrpcClient.GetUserById(ctx, &userPb.UserRequest{
		Id: int32(order.UserId),
	})
	if err != nil {
		fmt.Println(err.Error())
		// user deleted
		msg := "user deleted."
		orderUser := OrderUserResponse{UserName: &msg}
		orderResponse.User = orderUser
	} else {
		userId := int(user.GetId())
		orderUser := OrderUserResponse{
			UserId:   &userId,
			UserName: &user.Name,
		}
		orderResponse.User = orderUser
	}
	// get product info
	productResponses := make([]OrderProductsResponse, 0, len(order.Products))
	var orderPrice int
	for _, p := range order.Products {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		product, err := s.productGrpcClient.GetProductById(ctx, &productPb.ProductRequest{
			Id: int32(p.ProductId),
		})
		var productResponse OrderProductsResponse
		if err != nil {
			msg := "product deleted"
			productResponse = OrderProductsResponse{
				ProductName: &msg,
			}
		} else {
			productId := int(product.GetId())
			productAmount := p.ProductAmount
			productPrice := int(product.GetPrice())
			productResponse = OrderProductsResponse{
				&productId,
				&product.Name,
				&productPrice,
				&productAmount,
			}
			orderPrice += *productResponse.ProductPrice * *productResponse.ProductAmount
		}
		productResponses = append(productResponses, productResponse)
	}
	orderResponse.OrderPrice = orderPrice
	orderResponse.Products = productResponses
	return &orderResponse
}

func (s *service) Create(orderRequest OrderRequest) (*Order, error) {
	err := s.checkUserValid(orderRequest.UserId)

	if err != nil {
		// user not found
		return nil, err
	}

	err = s.checkProductsValid(orderRequest.OrderProduct)

	if err != nil {
		// product not found
		return nil, err
	}

	order := Order{
		UserId:   orderRequest.UserId,
		Products: orderRequest.OrderProduct,
	}
	newOrder, err := s.orderRepository.Create(order)
	return newOrder, err
}

func (s *service) Update(Id int, orderRequest OrderRequest) (*Order, error) {
	err := s.checkUserValid(orderRequest.UserId)

	if err != nil {
		// user not found
		return nil, err
	}

	err = s.checkProductsValid(orderRequest.OrderProduct)

	if err != nil {
		// product not found
		return nil, err
	}

	order, err := s.orderRepository.FindById(Id)

	if err != nil {
		return nil, err
	}

	order.UserId = orderRequest.UserId
	order.Products = orderRequest.OrderProduct
	newOrder, err := s.orderRepository.Update(*order)
	return newOrder, err
}

func (s *service) checkUserValid(userId int) error {
	// get user info
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.userGrpcClient.GetUserById(ctx, &userPb.UserRequest{
		Id: int32(userId),
	})

	if err != nil {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (s *service) checkProductsValid(orderProducts []OrderProduct) error {
	for _, p := range orderProducts {
		// get products info
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := s.productGrpcClient.GetProductById(ctx, &productPb.ProductRequest{
			Id: int32(p.ProductId),
		})

		if err != nil {
			return fmt.Errorf("product not found, id %v", p.ProductId)
		}
	}

	return nil
}

func (s *service) Delete(Id int) (*Order, error) {
	order, err := s.orderRepository.Delete(Id)
	return order, err
}

// func convertToOrdersResponse(rawOrders *[]order.Order) []order.OrderResponse {
// 	var ordersResponse []order.OrderResponse
// 	for _, u := range *rawOrders {
// 		orderResponse := convertToOrderResponse(&u)
// 		ordersResponse = append(ordersResponse, orderResponse)
// 	}
// 	return ordersResponse
// }

// func convertToOrderResponse(rawOrder *order.Order) order.OrderResponse {
// 	return order.OrderResponse{
// 		Id: rawOrder.Id,
// 	}
// }
