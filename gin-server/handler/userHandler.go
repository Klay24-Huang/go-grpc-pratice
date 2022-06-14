package handler

import (
	"context"
	"fmt"
	pb "grpc-practice/grpc-gateway/user"
	"grpc-practice/user-service/user"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type userHandler struct {
	userClient pb.UserClient
}

func NewUserHandler(userClient pb.UserClient) *userHandler {
	return &userHandler{userClient}
}

func (h *userHandler) GetUserList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	usersReply, err := h.userClient.GetUserList(ctx, &pb.Empty{})

	if err != nil {
		errMsg := fmt.Sprintf("get user list error, %v", err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	var usersResponse []user.UserResponse
	for _, u := range usersReply.GetUsers() {
		userResponse := user.UserResponse{
			Id:      int(u.GetId()),
			Account: u.GetAccount(),
			Name:    u.GetName(),
			Phone:   u.GetPhone(),
		}
		usersResponse = append(usersResponse, userResponse)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": usersResponse,
	})
}

func (h *userHandler) GetUserById(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	userReply, err := h.userClient.GetUserById(ctx, &pb.UserRequest{Id: int32(id)})
	if err != nil {
		errMsg := fmt.Sprintf("get user by id error, id: %v, %v", id, err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.UserResponse{
			Id:      int(userReply.GetId()),
			Account: userReply.GetAccount(),
			Name:    userReply.GetName(),
			Phone:   userReply.GetPhone(),
		},
	})
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var userRequest user.UserRequest

	err := c.ShouldBindJSON(&userRequest)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errMessage := fmt.Sprintf("Error on filled %s, condition: %s", e.Field(), e.ActualTag())
			c.JSON(http.StatusBadRequest, errMessage)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := h.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:    userRequest.Name,
		Account: userRequest.Account,
		Phone:   userRequest.Phone,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	var updateUserRequest user.UserRequest

	err := c.ShouldBindJSON(&updateUserRequest)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errMessage := fmt.Sprintf("Error on filled %s, condition: %s", e.Field(), e.ActualTag())
			c.JSON(http.StatusBadRequest, errMessage)
			return
		}
	}

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	updateUser, err := h.userClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:      int32(id),
		Name:    updateUserRequest.Name,
		Phone:   updateUserRequest.Phone,
		Account: updateUserRequest.Account,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updateUser,
	})
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := h.userClient.DeleteUser(ctx, &pb.UserRequest{Id: int32(id)})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Delete user success",
	})

}
