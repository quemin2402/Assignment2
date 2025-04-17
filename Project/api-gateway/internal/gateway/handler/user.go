package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userpb "github.com/quemin2402/user-service/proto"
)

type UserH struct {
	Cli userpb.UserServiceClient
}

func NewUser(cli userpb.UserServiceClient) *UserH {
	return &UserH{Cli: cli}
}

func (h *UserH) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Cli.RegisterUser(c, &userpb.UserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
