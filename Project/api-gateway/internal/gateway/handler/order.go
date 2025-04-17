package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderpb "github.com/quemin2402/order-service/proto"
)

type OrderH struct {
	Cli orderpb.OrderServiceClient
}

func NewOrder(cli orderpb.OrderServiceClient) *OrderH {
	return &OrderH{Cli: cli}
}

func (h *OrderH) Create(c *gin.Context) {
	var req struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Items  []struct {
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		} `json:"items"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []*orderpb.OrderItem
	for _, i := range req.Items {
		items = append(items, &orderpb.OrderItem{ProductId: i.ProductID, Quantity: i.Quantity})
	}

	_, err := h.Cli.CreateOrder(c, &orderpb.CreateOrderRequest{
		Order: &orderpb.Order{Id: req.ID, Status: req.Status, Items: items},
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
