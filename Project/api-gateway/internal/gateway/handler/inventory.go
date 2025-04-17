package handler

import (
	"github.com/gin-gonic/gin"
	invpb "github.com/quemin2402/inventory-service/proto"
	"net/http"
)

type Inv struct{ cli invpb.InventoryServiceClient }

func NewInv(c invpb.InventoryServiceClient) *Inv { return &Inv{c} }

func (h *Inv) Create(c *gin.Context) {
	var req struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Category string  `json:"category"`
		Price    float64 `json:"price"`
		Stock    int32   `json:"stock"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	_, err := h.cli.CreateProduct(c,
		&invpb.CreateProductRequest{Product: &invpb.Product{Id: req.ID, Name: req.Name, Category: req.Category, Price: req.Price, Stock: req.Stock}})
	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Inv) List(c *gin.Context) {
	stream, err := h.cli.ListProducts(c, &invpb.ListProductsRequest{})
	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var out []interface{}
	for {
		p, err := stream.Recv()
		if err != nil {
			break
		}
		out = append(out, p)
	}
	c.JSON(http.StatusOK, out)
}
