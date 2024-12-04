package controller

import (
	"encoding/json"
	"strconv"

	"go-api-mono/internal/app/user/model"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/auth"
	"go-api-mono/internal/pkg/core"
	"go-api-mono/internal/pkg/errors"
)

// UserController 用户控制器
type UserController struct {
	service *service.UserService
	jwt     *auth.JWT
}

// NewUserController 创建用户控制器
func NewUserController(service *service.UserService, jwt *auth.JWT) *UserController {
	return &UserController{
		service: service,
		jwt:     jwt,
	}
}

// Register 注册用户
func (c *UserController) Register(ctx *core.Context) {
	var req model.RegisterRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid request body"))
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.service.RegisterUser(ctx, user); err != nil {
		ctx.Response.Error(err)
		return
	}

	ctx.Response.Success(model.RegisterResponse{User: user})
}

// Login 用户登录
func (c *UserController) Login(ctx *core.Context) {
	var req model.LoginRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid request body"))
		return
	}

	user, err := c.service.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		ctx.Response.Error(err)
		return
	}

	token, err := c.jwt.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeInternal, "failed to generate token"))
		return
	}

	ctx.Response.Success(model.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(c.jwt.Config.ExpirationTime.Seconds()),
		User:      user,
	})
}

// List 获取用户列表
func (c *UserController) List(ctx *core.Context) {
	page, _ := strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(ctx.Request.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := c.service.ListUsers(ctx, page, pageSize)
	if err != nil {
		ctx.Response.Error(err)
		return
	}

	ctx.Response.Success(map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// Get 获取用户详情
func (c *UserController) Get(ctx *core.Context) {
	id, err := strconv.ParseUint(ctx.Request.PathValue("id"), 10, 32)
	if err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid user id"))
		return
	}

	user, err := c.service.GetUser(ctx, uint(id))
	if err != nil {
		ctx.Response.Error(err)
		return
	}

	ctx.Response.Success(user)
}

// Update 更新用户
func (c *UserController) Update(ctx *core.Context) {
	id, err := strconv.ParseUint(ctx.Request.PathValue("id"), 10, 32)
	if err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid user id"))
		return
	}

	var user model.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid request body"))
		return
	}

	user.ID = uint(id)
	if err := c.service.UpdateUser(ctx, &user); err != nil {
		ctx.Response.Error(err)
		return
	}

	ctx.Response.Success(user)
}

// Delete 删除用户
func (c *UserController) Delete(ctx *core.Context) {
	id, err := strconv.ParseUint(ctx.Request.PathValue("id"), 10, 32)
	if err != nil {
		ctx.Response.Error(errors.New(errors.ErrCodeBadRequest, "invalid user id"))
		return
	}

	if err := c.service.DeleteUser(ctx, uint(id)); err != nil {
		ctx.Response.Error(err)
		return
	}

	ctx.Response.NoContent()
}
