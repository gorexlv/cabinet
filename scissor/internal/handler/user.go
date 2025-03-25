package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorexlv/cabinet/scissor/internal/domain"
	"github.com/gorexlv/cabinet/scissor/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/users").To(h.Create).
		Doc("创建用户").
		Reads(domain.CreateUserRequest{}).
		Returns(200, "OK", domain.User{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.POST("/users/login").To(h.Login).
		Doc("用户登录").
		Reads(domain.LoginRequest{}).
		Returns(200, "OK", domain.LoginResponse{}).
		Returns(401, "Unauthorized", nil))

	ws.Route(ws.POST("/users/wx-login").To(h.WxLogin).
		Doc("微信登录").
		Reads(domain.WxLoginRequest{}).
		Returns(200, "OK", domain.LoginResponse{}).
		Returns(400, "Bad Request", nil))
}

func (h *UserHandler) Create(req *restful.Request, resp *restful.Response) {
	var createReq domain.CreateUserRequest
	if err := req.ReadEntity(&createReq); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Create(req.Request.Context(), &createReq)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	resp.WriteEntity(user)
}

func (h *UserHandler) Login(req *restful.Request, resp *restful.Response) {
	var loginReq domain.LoginRequest
	if err := req.ReadEntity(&loginReq); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	loginResp, err := h.userService.Login(req.Request.Context(), &loginReq)
	if err != nil {
		resp.WriteError(http.StatusUnauthorized, err)
		return
	}

	resp.WriteEntity(loginResp)
}

func (h *UserHandler) WxLogin(req *restful.Request, resp *restful.Response) {
	var wxLoginReq domain.WxLoginRequest
	if err := req.ReadEntity(&wxLoginReq); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	loginResp, err := h.userService.WxLogin(req.Request.Context(), &wxLoginReq)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	resp.WriteEntity(loginResp)
}
