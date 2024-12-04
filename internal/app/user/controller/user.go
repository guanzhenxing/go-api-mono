package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/response"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register registers user-related routes
func (c *UserController) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", c.Create)
		users.GET("", c.List)
		users.GET("/:id", c.Get)
		users.PUT("/:id", c.Update)
		users.DELETE("/:id", c.Delete)
	}
}

// Create handles user creation
func (c *UserController) Create(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	if err := c.userService.Register(ctx, &user); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
		"data": gin.H{
			"user": user,
		},
	})
}

// Get handles getting a user by ID
func (c *UserController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := c.userService.Get(ctx, uint(id))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
		"data":    user,
	})
}

// List handles listing all users
func (c *UserController) List(ctx *gin.Context) {
	users, err := c.userService.List(ctx)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
		"data": gin.H{
			"page":  1,
			"size":  10,
			"total": len(users),
			"users": users,
		},
	})
}

// Update handles updating a user
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}
	user.ID = uint(id)

	if err := c.userService.Update(ctx, &user); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
		"data":    user,
	})
}

// Delete handles deleting a user
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err)
		return
	}

	if err := c.userService.Delete(ctx, uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
