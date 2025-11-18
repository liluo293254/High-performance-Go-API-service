package user

import (
	"bgame/internal/service"
	"bgame/internal/util"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// RegAndLogin 用户注册和登录合并
// @Summary      用户注册和登录合并
// @Description  用户注册和登录合并接口
// @Tags         用户接口
// @Accept       json
// @Produce      json
// @Param        request body service.RegAndLoginRequest true "注册和登录请求"
// @Success      200  {object}  util.Response{data=service.RegAndLoginResponse}
// @Failure      400  {object}  util.Response
// @Router       /api/user/regAndLogin [post]
func (h *UserHandler) RegAndLogin(c *gin.Context) {
	var req service.RegAndLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, "参数错误: "+err.Error())
		return
	}
	resp, err := h.userService.RegAndLogin(&req)
	if err != nil {
		util.Error(c, err.Error())
		return
	}
	util.Success(c, resp)
}

// GetUserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  获取用户信息接口
// @Tags         用户接口
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  util.Response{data=model.User}
// @Failure      400  {object}  util.Response
// @Router       /api/user/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		util.Unauthorized(c, "未获取到用户信息")
		return
	}
	user, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		util.Error(c, err.Error())
		return
	}
	util.Success(c, user)
}