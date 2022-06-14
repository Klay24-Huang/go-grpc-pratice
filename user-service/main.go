package main

import (
	"context"
	"log"
	"net"

	pb "grpc-practice/grpc-gateway/user"
	"grpc-practice/user-service/user"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

// Server ...
type Server struct {
	pb.UnimplementedUserServer
}

var userService user.Service

// DI repository and service
func InjectionService() {
	userRepository := user.NewRepository()
	userService = user.NewService(userRepository)
}

func (s *Server) GetUserById(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	//log.Printf("Received: %v", in.GetId())
	id := in.GetId()
	user, err := userService.FindById(int(id))

	if err != nil {
		return nil, err
	}

	return &pb.UserReply{Id: int32(user.Id), Account: user.Account, Name: user.Account, Phone: user.Phone}, nil
}

func (s *Server) GetUserList(ctx context.Context, in *pb.Empty) (*pb.UsersReply, error) {
	users, err := userService.FindAll()

	if err != nil {
		return nil, err
	}

	var usersReply []*pb.UserReply
	for _, u := range *users {
		userReply := pb.UserReply{Id: int32(u.Id), Account: u.Account, Name: u.Name, Phone: u.Phone}
		usersReply = append(usersReply, &userReply)
	}
	return &pb.UsersReply{Users: usersReply}, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UserReply, error) {
	id := int(in.GetId())
	request := user.UserRequest{
		Name:    in.Name,
		Account: in.Account,
		Phone:   in.Phone,
	}
	newUser, err := userService.Update(id, request)

	if err != nil {
		return nil, err
	}
	return &pb.UserReply{Id: int32(newUser.Id), Account: newUser.Account, Name: newUser.Name, Phone: newUser.Phone}, nil
}

func (s *Server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserReply, error) {
	userRequest := user.UserRequest{
		Name:    in.GetName(),
		Account: in.GetAccount(),
		Phone:   in.GetPhone(),
	}

	user, err := userService.Create(userRequest)

	if err != nil {
		return nil, err
	}

	return &pb.UserReply{Id: int32(user.Id), Account: user.Account, Name: user.Name, Phone: user.Phone}, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	id := in.GetId()

	user, err := userService.Delete(int(id))

	if err != nil {
		return nil, err
	}

	return &pb.UserReply{Id: int32(user.Id), Account: user.Account, Name: user.Name, Phone: user.Phone}, nil
}

func main() {
	// repository injection
	InjectionService()

	// run grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
