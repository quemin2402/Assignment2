package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/quemin2402/api-gateway/internal/gateway/client"
	"github.com/quemin2402/api-gateway/internal/gateway/handler"
)

func NewRouter() *gin.Engine {
	invCli, _, _ := client.NewInventory("localhost:9001")
	ordCli, _, _ := client.NewOrder("localhost:9002")
	usrCli, _, _ := client.NewUser("localhost:9003")

	r := gin.Default()

	// INVENTORY
	inv := handler.NewInv(invCli)
	r.POST("/products", inv.Create)
	r.GET("/products", inv.List)

	// ORDER
	ord := handler.NewOrder(ordCli)
	r.POST("/orders", ord.Create)

	// USER
	usr := handler.NewUser(usrCli)
	r.POST("/users/register", usr.Register)

	return r
}
