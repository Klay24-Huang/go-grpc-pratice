package handler

import (
	"context"
	"fmt"
	pb "grpc-practice/grpc-gateway/order"
	"grpc-practice/order-service/order"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type orderHandler struct {
	orderClient pb.OrderClient
}

func NewOrderHandler(orderClient pb.OrderClient) *orderHandler {
	return &orderHandler{orderClient}
}

// @Summary 取得order列表
// @Tags order
// @product application/json
// @Success 200 {array} order.OrderResponse
// @Router /orders [get]
func (h *orderHandler) GetOrderList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ordersReply, err := h.orderClient.GetOrderList(ctx, &pb.Empty{})

	if err != nil {
		errMsg := fmt.Sprintf("get order list error, %v", err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	var ordersResponse []order.OrderResponse
	for _, o := range ordersReply.GetOrders() {
		userId := int(o.User.UserId)
		orderResponse := order.OrderResponse{
			Id: int(o.GetId()),
			User: order.OrderUserResponse{
				UserId:   &userId,
				UserName: &o.User.UserName,
			},
			OrderPrice: int(o.GetOrderPrice()),
		}
		// product
		var productsResponse []order.OrderProductsResponse
		for _, p := range o.GetProducts() {
			productId := int(p.GetProductId())
			productPrice := int(p.GetProductPrice())
			productAmount := int(p.GetProductAmount())

			productResponse := order.OrderProductsResponse{
				ProductId:     &productId,
				ProductName:   &p.ProductName,
				ProductPrice:  &productPrice,
				ProductAmount: &productAmount,
			}
			productsResponse = append(productsResponse, productResponse)
		}
		orderResponse.Products = productsResponse
		ordersResponse = append(ordersResponse, orderResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ordersResponse,
	})
}

// @Summary 取得指定id order資訊
// @Tags order
// @product application/json
// @param order_id path int true "order id"
// @Success 200 {object} order.OrderResponse
// @Router /orders/{order_id} [get]
func (h *orderHandler) GetOrderById(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	orderReply, err := h.orderClient.GetOrderById(ctx, &pb.OrderRequest{Id: int32(id)})
	if err != nil {
		errMsg := fmt.Sprintf("get order by id error, id: %v, %v", id, err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	user := orderReply.GetUser()
	userId := int(user.GetUserId())
	userName := user.GetUserName()
	userResponse := order.OrderUserResponse{
		UserId:   &userId,
		UserName: &userName,
	}

	products := orderReply.GetProducts()
	var productsResponse []order.OrderProductsResponse
	for _, p := range products {
		productId := int(p.GetProductId())
		productName := p.GetProductName()
		productPrice := int(p.GetProductPrice())
		productAmount := int(p.GetProductAmount())
		productResponse := order.OrderProductsResponse{
			ProductId:     &productId,
			ProductName:   &productName,
			ProductPrice:  &productPrice,
			ProductAmount: &productAmount,
		}
		productsResponse = append(productsResponse, productResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order.OrderResponse{
			Id:         int(orderReply.GetId()),
			User:       userResponse,
			Products:   productsResponse,
			OrderPrice: int(orderReply.GetOrderPrice()),
		},
	})
}

// @Summary 新增order
// @Tags order
// @product application/json
// @param order body order.OrderRequest true "order info"
// @Success 200 {object} pb.OrderReply
// @Router /orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	var orderRequest order.OrderRequest

	err := c.ShouldBindJSON(&orderRequest)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errMessage := fmt.Sprintf("Error on filled %s, condition: %s", e.Field(), e.ActualTag())
			c.JSON(http.StatusBadRequest, errMessage)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var orderProducts []*pb.UpdateOrderProductRequest
	for _, p := range orderRequest.OrderProduct {
		orderProduct := pb.UpdateOrderProductRequest{
			ProductId:     int32(p.ProductId),
			ProductAmount: int32(p.ProductAmount),
		}
		orderProducts = append(orderProducts, &orderProduct)
	}

	order, err := h.orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId:   int32(orderRequest.UserId),
		Products: orderProducts,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}

// @Summary 更新order資訊
// @Tags order
// @product application/json
// @param id path int true "id"
// @param order body order.OrderRequest true "order info"
// @Success 200 {object} pb.OrderReply
// @Router /orders [put]
func (h *orderHandler) UpdateOrder(c *gin.Context) {
	var updateOrderRequest order.OrderRequest

	err := c.ShouldBindJSON(&updateOrderRequest)

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

	var updateProductRequests []*pb.UpdateOrderProductRequest
	for _, p := range updateOrderRequest.OrderProduct {
		updateProductRequest := pb.UpdateOrderProductRequest{
			ProductId:     int32(p.ProductId),
			ProductAmount: int32(p.ProductAmount),
		}
		updateProductRequests = append(updateProductRequests, &updateProductRequest)
	}
	updateOrder, err := h.orderClient.UpdateOrder(ctx, &pb.UpdateOrderRequest{
		Id:      int32(id),
		UserId:  int32(updateOrderRequest.UserId),
		Product: updateProductRequests,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": updateOrder,
	})
}

// @Summary 刪除指定id order
// @Tags order
// @product application/json
// @param order_id path int true "order id"
// @Success 200
// @Router /orders/{order_id} [delete]
func (h *orderHandler) DeleteOrder(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	order, err := h.orderClient.DeleteOrder(ctx, &pb.OrderRequest{Id: int32(id)})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    order,
		"message": "Delete order success",
	})
}
