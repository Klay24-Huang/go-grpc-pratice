package handler

import (
	"context"
	"fmt"
	pb "grpc-practice/grpc-gateway/product"
	"grpc-practice/product-service/product"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type productHandler struct {
	productClient pb.ProductClient
}

func NewProductHandler(productClient pb.ProductClient) *productHandler {
	return &productHandler{productClient}
}

func (h *productHandler) GetProductList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	productsReply, err := h.productClient.GetProductList(ctx, &pb.Empty{})

	if err != nil {
		errMsg := fmt.Sprintf("get product list error, %v", err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	var productsResponse []product.ProductResponse
	for _, p := range productsReply.GetProducts() {
		productResponse := product.ProductResponse{
			Id:    int(p.GetId()),
			Name:  p.GetName(),
			Price: int(p.Price),
		}
		productsResponse = append(productsResponse, productResponse)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": productsResponse,
	})
}

func (h *productHandler) GetProductById(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	productReply, err := h.productClient.GetProductById(ctx, &pb.ProductRequest{Id: int32(id)})
	if err != nil {
		errMsg := fmt.Sprintf("get product by id error, id: %v, %v", id, err)
		//log.Fatalf(errMsg))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product.ProductResponse{
			Id:    int(productReply.GetId()),
			Name:  productReply.GetName(),
			Price: int(productReply.GetPrice()),
		},
	})
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var productRequest product.ProductRequest

	err := c.ShouldBindJSON(&productRequest)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errMessage := fmt.Sprintf("Error on filled %s, condition: %s", e.Field(), e.ActualTag())
			c.JSON(http.StatusBadRequest, errMessage)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	product, err := h.productClient.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:  productRequest.Name,
		Price: int32(productRequest.Price),
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	var updateProductRequest product.ProductRequest

	err := c.ShouldBindJSON(&updateProductRequest)

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
	updateProduct, err := h.productClient.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:    int32(id),
		Name:  updateProductRequest.Name,
		Price: int32(updateProductRequest.Price),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": updateProduct,
	})
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	product, err := h.productClient.DeleteProduct(ctx, &pb.ProductRequest{Id: int32(id)})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    product,
		"message": "Delete product success",
	})

}
