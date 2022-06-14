package main

import (
	"context"
	"log"
	"net"

	pb "grpc-practice/grpc-gateway/order"
	productPb "grpc-practice/grpc-gateway/product"
	userPb "grpc-practice/grpc-gateway/user"
	"grpc-practice/order-service/order"

	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	pb.UnimplementedOrderServer
}

var orderService order.Service

func (s *Server) GetOrderById(ctx context.Context, in *pb.OrderRequest) (*pb.OrderDetailReply, error) {
	//log.Printf("Received: %v", in.GetId())
	id := in.GetId()
	order, err := orderService.FindById(int(id))

	if err != nil {
		return nil, err
	}

	user := order.User
	products := order.Products
	var productsReply []*pb.OrderProductDetailReply
	for _, product := range products {
		productReply := pb.OrderProductDetailReply{
			ProductId:     int32(*product.ProductId),
			ProductName:   *product.ProductName,
			ProductPrice:  int32(*product.ProductPrice),
			ProductAmount: int32(*product.ProductAmount),
		}
		productsReply = append(productsReply, &productReply)
	}

	return &pb.OrderDetailReply{
		Id: int32(order.Id),
		User: &pb.OrderUserReply{
			UserId:   int32(*user.UserId),
			UserName: *user.UserName,
		},
		Products:   productsReply,
		OrderPrice: int32(order.OrderPrice),
	}, nil
}

func (s *Server) GetOrderList(ctx context.Context, in *pb.Empty) (*pb.OrdersReply, error) {
	orders, err := orderService.FindAll()

	if err != nil {
		return nil, err
	}

	var ordersReply []*pb.OrderDetailReply
	for _, o := range *orders {
		products := o.Products
		var productsReply []*pb.OrderProductDetailReply
		for _, product := range products {
			productReply := pb.OrderProductDetailReply{
				ProductId:     int32(*product.ProductId),
				ProductName:   *product.ProductName,
				ProductPrice:  int32(*product.ProductPrice),
				ProductAmount: int32(*product.ProductAmount),
			}
			productsReply = append(productsReply, &productReply)
		}

		orderReply := pb.OrderDetailReply{
			Id: int32(o.Id),
			User: &pb.OrderUserReply{
				UserId:   int32(*o.User.UserId),
				UserName: *o.User.UserName,
			},
			Products:   productsReply,
			OrderPrice: int32(o.OrderPrice),
		}
		ordersReply = append(ordersReply, &orderReply)
	}
	return &pb.OrdersReply{Orders: ordersReply}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, in *pb.UpdateOrderRequest) (*pb.OrderReply, error) {
	id := int(in.GetId())
	var orderProducts []order.OrderProduct
	for _, product := range in.GetProduct() {
		orderProduct := order.OrderProduct{
			ProductId:     int(product.GetProductId()),
			ProductAmount: int(product.GetProductAmount()),
		}
		orderProducts = append(orderProducts, orderProduct)
	}

	request := order.OrderRequest{
		UserId:       int(in.GetUserId()),
		OrderProduct: orderProducts,
	}
	newOrder, err := orderService.Update(id, request)
	if err != nil {
		return nil, err
	}

	return &pb.OrderReply{
		Id:       int32(newOrder.Id),
		UserId:   int32(newOrder.UserId),
		Products: convertOrderProductReply(&newOrder.Products),
	}, nil
}

func convertOrderProductReply(orderProducts *[]order.OrderProduct) []*pb.OrderProductReply {
	var orderProductsReply []*pb.OrderProductReply
	for _, orderProduct := range *orderProducts {
		orderProductReply := pb.OrderProductReply{
			ProductId:     int32(orderProduct.ProductId),
			ProductAmount: int32(orderProduct.ProductAmount),
		}
		orderProductsReply = append(orderProductsReply, &orderProductReply)
	}
	return orderProductsReply
}

func (s *Server) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderReply, error) {
	var orderProducts []order.OrderProduct
	for _, rawProduct := range in.GetProducts() {
		orderProduct := order.OrderProduct{
			ProductId:     int(rawProduct.GetProductId()),
			ProductAmount: int(rawProduct.GetProductAmount()),
		}
		orderProducts = append(orderProducts, orderProduct)
	}

	orderRequest := order.OrderRequest{
		UserId:       int(in.GetUserId()),
		OrderProduct: orderProducts,
	}

	order, err := orderService.Create(orderRequest)

	if err != nil {
		return nil, err
	}

	return &pb.OrderReply{
		Id:       int32(order.Id),
		UserId:   int32(order.UserId),
		Products: convertOrderProductReply(&order.Products),
	}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, in *pb.OrderRequest) (*pb.OrderReply, error) {
	id := in.GetId()

	order, err := orderService.Delete(int(id))

	if err != nil {
		return nil, err
	}

	return &pb.OrderReply{
		Id:       int32(order.Id),
		UserId:   int32(order.UserId),
		Products: convertOrderProductReply(&order.Products),
	}, nil
}

const (
	port                  = ":50054"
	userServiceAddress    = "localhost:50052"
	productServiceAddress = "localhost:50053"
)

func main() {
	// Injection
	userServiceConn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect user service: %v", err)
	}
	defer userServiceConn.Close()
	userGrpcClient := userPb.NewUserClient(userServiceConn)

	productServiceConn, err := grpc.Dial(productServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect product service: %v", err)
	}
	defer productServiceConn.Close()
	productGrpcClient := productPb.NewProductClient(productServiceConn)

	orderRepository := order.NewRepository()
	orderService = order.NewService(orderRepository, userGrpcClient, productGrpcClient)

	// run grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
