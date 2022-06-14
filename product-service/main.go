package main

import (
	"context"
	"log"
	"net"

	pb "grpc-practice/grpc-gateway/product"
	"grpc-practice/product-service/product"

	"google.golang.org/grpc"
)

const (
	port = ":50053"
)

// Server ...
type Server struct {
	pb.UnimplementedProductServer
}

var productService product.Service

// DI repository and service
func InjectionService() {
	productRepository := product.NewRepository()
	productService = product.NewService(productRepository)
}

func (s *Server) GetProductById(ctx context.Context, in *pb.ProductRequest) (*pb.ProductReply, error) {
	//log.Printf("Received: %v", in.GetId())
	id := in.GetId()
	product, err := productService.FindById(int(id))

	if err != nil {
		return nil, err
	}

	return &pb.ProductReply{Id: int32(product.Id), Name: product.Name, Price: int32(product.Price)}, nil
}

func (s *Server) GetProductList(ctx context.Context, in *pb.Empty) (*pb.ProductsReply, error) {
	products, err := productService.FindAll()

	if err != nil {
		return nil, err
	}

	var productsReply []*pb.ProductReply
	for _, p := range *products {
		productReply := pb.ProductReply{Id: int32(p.Id), Name: p.Name, Price: int32(p.Price)}
		productsReply = append(productsReply, &productReply)
	}
	return &pb.ProductsReply{Products: productsReply}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.ProductReply, error) {
	id := int(in.GetId())
	request := product.ProductRequest{
		Name:  in.Name,
		Price: int(in.Price),
	}
	newProduct, err := productService.Update(id, request)

	if err != nil {
		return nil, err
	}
	return &pb.ProductReply{Id: int32(newProduct.Id), Name: newProduct.Name, Price: int32(newProduct.Price)}, nil
}

func (s *Server) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.ProductReply, error) {
	productRequest := product.ProductRequest{
		Name:  in.GetName(),
		Price: int(in.GetPrice()),
	}

	product, err := productService.Create(productRequest)

	if err != nil {
		return nil, err
	}

	return &pb.ProductReply{Id: int32(product.Id), Name: product.Name, Price: int32(product.Price)}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, in *pb.ProductRequest) (*pb.ProductReply, error) {
	id := in.GetId()

	product, err := productService.Delete(int(id))

	if err != nil {
		return nil, err
	}

	return &pb.ProductReply{Id: int32(product.Id), Name: product.Name, Price: int32(product.Price)}, nil
}

func main() {
	InjectionService()
	// run grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
