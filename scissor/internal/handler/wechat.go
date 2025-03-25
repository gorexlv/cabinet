package handler

import (
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/gorexlv/cabinet/scissor/internal/domain"
	"github.com/gorexlv/cabinet/scissor/internal/service"
	"github.com/gorexlv/cabinet/scissor/pkg/wechat"
)

// WechatHandler 处理微信相关的请求
type WechatHandler struct {
	wxClient *wechat.Client
	userSvc  *service.UserService
}

// NewWechatHandler 创建微信处理器
func NewWechatHandler(wxClient *wechat.Client, userSvc *service.UserService) *WechatHandler {
	return &WechatHandler{
		wxClient: wxClient,
		userSvc:  userSvc,
	}
}

// Register 注册路由
func (h *WechatHandler) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/wechat/qrcode").To(h.GetQRCode).
		Doc("获取微信登录二维码").
		Returns(200, "OK", nil))

	ws.Route(ws.GET("/wechat/callback").To(h.HandleCallback).
		Doc("处理微信登录回调").
		Param(ws.QueryParameter("code", "微信授权码").Required(true)).
		Returns(200, "OK", nil))

	ws.Route(ws.GET("/wechat/check").To(h.CheckLoginStatus).
		Doc("检查登录状态").
		Param(ws.QueryParameter("code", "微信授权码").Required(true)).
		Returns(200, "OK", nil))
}

// GetQRCode 获取微信登录二维码
func (h *WechatHandler) GetQRCode(req *restful.Request, resp *restful.Response) {
	qrcodeURL, err := h.wxClient.GetQRCode()
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("生成二维码失败: %v", err),
		}, "application/json")
		return
	}

	resp.WriteHeaderAndJson(http.StatusOK, map[string]string{
		"qrcode_url": qrcodeURL,
	}, "application/json")
}

// HandleCallback 处理微信登录回调
func (h *WechatHandler) HandleCallback(req *restful.Request, resp *restful.Response) {
	code := req.QueryParameter("code")
	if code == "" {
		resp.WriteHeaderAndJson(http.StatusBadRequest, map[string]string{
			"error": "缺少code参数",
		}, "application/json")
		return
	}

	// 获取微信用户信息
	wxResp, err := h.wxClient.CheckLogin(code)
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("获取微信token失败: %v", err),
		}, "application/json")
		return
	}

	userInfo, err := h.wxClient.GetUserInfo(wxResp.AccessToken, wxResp.OpenID)
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("获取用户信息失败: %v", err),
		}, "application/json")
		return
	}

	// 查找或创建用户
	user, err := h.userSvc.WxLogin(req.Request.Context(), &domain.WxLoginRequest{
		Code:     code,
		OpenID:   userInfo.OpenID,
		Nickname: userInfo.Nickname,
	})
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("处理用户信息失败: %v", err),
		}, "application/json")
		return
	}

	resp.WriteHeaderAndJson(http.StatusOK, user, "application/json")
}

// CheckLoginStatus 检查登录状态
func (h *WechatHandler) CheckLoginStatus(req *restful.Request, resp *restful.Response) {
	code := req.QueryParameter("code")
	if code == "" {
		resp.WriteHeaderAndJson(http.StatusBadRequest, map[string]string{
			"error": "缺少code参数",
		}, "application/json")
		return
	}

	wxResp, err := h.wxClient.CheckLogin(code)
	if err != nil {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("检查登录状态失败: %v", err),
		}, "application/json")
		return
	}

	resp.WriteHeaderAndJson(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   wxResp,
	}, "application/json")
}
