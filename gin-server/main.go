// main.go
package main

import (
	"log"
	"net/http"

	"grpc-practice/gin-server/handler"
	orderPb "grpc-practice/grpc-gateway/order"
	productPb "grpc-practice/grpc-gateway/product"
	userPb "grpc-practice/grpc-gateway/user"

	_ "grpc-practice/gin-server/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userServiceAddress    = "user-service:50052"
	productServiceAddress = "product-service:50053"
	orderServiceAddress   = "order-service:50054"
)

// @title Go gRPC simple demo
// @version 1.0
// @description Golang 微服務範例
// @description 更多說明參考 https://github.com/Klay24-Huang/go-grpc-pratice
func main() {
	router := gin.Default()

	// user api
	// Set up a connection to the server.
	userServiceConn, err := grpc.Dial(userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect user service: %v", err)
	}
	defer userServiceConn.Close()
	userGrpcClient := userPb.NewUserClient(userServiceConn)
	userHandler := handler.NewUserHandler(userGrpcClient)

	router.GET("/users", userHandler.GetUserList)
	router.GET("/users/:id", userHandler.GetUserById)
	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	//router.DELETE("/users/:id", userHandler.DeleteUser)

	// // product api
	productServiceConn, err := grpc.Dial(productServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect product service: %v", err)
	}
	defer productServiceConn.Close()
	productGrpcClient := productPb.NewProductClient(productServiceConn)
	productHandler := handler.NewProductHandler(productGrpcClient)

	router.GET("/products", productHandler.GetProductList)
	router.GET("/products/:id", productHandler.GetProductById)
	router.POST("/products", productHandler.CreateProduct)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	//router.DELETE("/products/:id", productHandler.DeleteProduct)

	// // order api
	orderServiceConn, err := grpc.Dial(orderServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect order service: %v", err)
	}
	defer orderServiceConn.Close()
	orderGrpcClient := orderPb.NewOrderClient(orderServiceConn)
	orderHandler := handler.NewOrderHandler(orderGrpcClient)

	router.GET("/orders", orderHandler.GetOrderList)
	router.GET("/orders/:id", orderHandler.GetOrderById)
	router.POST("/orders", orderHandler.CreateOrder)
	router.PUT("/orders/:id", orderHandler.UpdateOrder)
	router.DELETE("/orders/:id", orderHandler.DeleteOrder)

	// test api
	router.GET("/health", test)

	// swagger
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(":1231")
}

func test(c *gin.Context) {
	str := []byte("ok")                      //因為網頁傳輸沒有string的概念，都是要轉成byte字節方式進行傳輸
	c.Data(http.StatusOK, "text/plain", str) // 指定contentType為 text/plain，就是傳輸格式為純文字
}

// create swagger json
// swag init -d  .,../user-service,../product-service,../order-service,../grpc-gateway/user,../grpc-gateway/product,../grpc-gateway/order
